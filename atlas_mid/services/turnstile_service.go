package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web"
)

const (
	turnstileValidationModeTest   = "test"
	turnstileValidationModeStrict = "strict"

	turnstileSiteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	turnstileTestSecretAlwaysPass = "1x0000000000000000000000000000000AA"

	turnstileTestSecretAlwaysFail = "2x0000000000000000000000000000000AA"

	turnstileTestSecretSpent = "3x0000000000000000000000000000000AA"
)

type turnstileConfig struct {
	ValidationMode  string
	Secret          string
	ExpectedHost    string
	ExpectedAction  string
	ClearanceSecret string
}

var ErrTurnstileUnavailable = errors.New(
	"cloudflare turnstile no disponible",
)

type turnstileVerifyRequest struct {
	Secret string `json:"secret"`

	Response string `json:"response"`

	RemoteIP string `json:"remoteip,omitempty"`
}

type TurnstileVerifyResponse struct {
	Success bool `json:"success"`

	ChallengeTS string `json:"challenge_ts,omitempty"`

	Hostname string `json:"hostname,omitempty"`

	Action string `json:"action,omitempty"`

	ErrorCodes []string `json:"error-codes,omitempty"`
}

var turnstileAvailability = struct {
	sync.RWMutex

	DegradedUntil time.Time
}{}

func getTurnstileConfig() turnstileConfig {
	return turnstileConfig{
		ValidationMode: strings.ToLower(
			strings.TrimSpace(
				web.AppConfig.DefaultString(
					"TurnstileValidationMode",
					"",
				),
			),
		),

		Secret: strings.TrimSpace(
			web.AppConfig.DefaultString(
				"TurnstileSecret",
				"",
			),
		),

		ExpectedHost: strings.TrimSpace(
			web.AppConfig.DefaultString(
				"TurnstileExpectedHostname",
				"",
			),
		),

		ExpectedAction: strings.TrimSpace(
			web.AppConfig.DefaultString(
				"TurnstileExpectedAction",
				"",
			),
		),

		ClearanceSecret: strings.TrimSpace(
			web.AppConfig.DefaultString(
				"TurnstileClearanceSecret",
				"",
			),
		),
	}
}

func ValidateTurnstileConfiguration() error {
	config := getTurnstileConfig()

	if config.ValidationMode == "" {
		return errors.New(
			"TurnstileValidationMode no configurado",
		)
	}

	if config.Secret == "" ||
		strings.HasPrefix(config.Secret, "${") {

		return errors.New(
			"TurnstileSecret no configurado",
		)
	}

	if config.ClearanceSecret == "" ||
		strings.HasPrefix(
			config.ClearanceSecret,
			"${",
		) {

		return errors.New(
			"TurnstileClearanceSecret no configurado",
		)
	}

	switch config.ValidationMode {

	case turnstileValidationModeTest:

		if !isTurnstileTestSecret(
			config.Secret,
		) {
			return errors.New(
				"modo Turnstile test requiere una secret oficial de pruebas",
			)
		}

	case turnstileValidationModeStrict:

		if isTurnstileTestSecret(
			config.Secret,
		) {
			return errors.New(
				"una secret Turnstile de pruebas no puede utilizarse en modo strict",
			)
		}

		if config.ExpectedHost == "" ||
			strings.HasPrefix(
				config.ExpectedHost,
				"${",
			) {

			return errors.New(
				"TurnstileExpectedHostname es obligatorio en modo strict",
			)
		}

		if config.ExpectedAction == "" ||
			strings.HasPrefix(
				config.ExpectedAction,
				"${",
			) {

			return errors.New(
				"TurnstileExpectedAction es obligatorio en modo strict",
			)
		}

	default:

		return fmt.Errorf(
			"TurnstileValidationMode inválido: %s",
			config.ValidationMode,
		)
	}

	return nil
}

func isTurnstileTestSecret(
	secret string,
) bool {

	switch secret {

	case turnstileTestSecretAlwaysPass,
		turnstileTestSecretAlwaysFail,
		turnstileTestSecretSpent:

		return true

	default:

		return false
	}
}

func ValidateTurnstile(
	token string,
	remoteIP string,
) (
	*TurnstileVerifyResponse,
	error,
) {

	token = strings.TrimSpace(
		token,
	)

	if token == "" {
		return nil,
			errors.New(
				"token Turnstile vacío",
			)
	}

	config :=
		getTurnstileConfig()

	requestBody :=
		turnstileVerifyRequest{
			Secret: config.Secret,

			Response: token,

			RemoteIP: remoteIP,
		}

	payload, err :=
		json.Marshal(
			requestBody,
		)

	if err != nil {
		return nil,
			fmt.Errorf(
				"error construyendo validación Turnstile: %w",
				err,
			)
	}

	client :=
		&http.Client{
			Timeout: 5 * time.Second,
		}

	response, err :=
		client.Post(
			turnstileSiteverifyURL,
			"application/json",
			bytes.NewReader(
				payload,
			),
		)

	if err != nil {

		markTurnstileUnavailable()

		return nil,
			fmt.Errorf(
				"%w: %v",
				ErrTurnstileUnavailable,
				err,
			)
	}

	defer response.Body.Close()

	if response.StatusCode >=
		http.StatusInternalServerError {

		markTurnstileUnavailable()

		return nil,
			fmt.Errorf(
				"%w: Siteverify respondió %s",
				ErrTurnstileUnavailable,
				response.Status,
			)
	}

	if response.StatusCode <
		http.StatusOK ||
		response.StatusCode >=
			http.StatusMultipleChoices {

		return nil,
			fmt.Errorf(
				"Siteverify respondió %s",
				response.Status,
			)
	}

	var result TurnstileVerifyResponse

	if err :=
		json.NewDecoder(
			response.Body,
		).Decode(
			&result,
		); err != nil {

		return nil,
			fmt.Errorf(
				"error procesando respuesta Siteverify: %w",
				err,
			)
	}

	markTurnstileAvailable()

	if !result.Success {
		return &result, nil
	}

	if config.ValidationMode ==
		turnstileValidationModeStrict {

		if result.Hostname !=
			config.ExpectedHost {

			result.Success = false

			result.ErrorCodes =
				append(
					result.ErrorCodes,
					"hostname-mismatch",
				)

			return &result, nil
		}

		if result.Action !=
			config.ExpectedAction {

			result.Success = false

			result.ErrorCodes =
				append(
					result.ErrorCodes,
					"action-mismatch",
				)

			return &result, nil
		}
	}

	return &result, nil
}

func GenerateTurnstileClearance(
	clientKey string,
) (string, error) {

	secret :=
		clearanceSecret()

	if secret == "" {
		return "",
			errors.New(
				"TurnstileClearanceSecret no configurado",
			)
	}

	ttl :=
		time.Duration(
			configInt(
				"TurnstileClearanceTTLSeconds",
				600,
			),
		) * time.Second

	expiresAt :=
		time.Now().
			Add(ttl).
			Unix()

	payload :=
		fmt.Sprintf(
			"%s.%d",
			hashClientKey(
				clientKey,
			),
			expiresAt,
		)

	signature :=
		signClearance(
			payload,
			secret,
		)

	return base64.
			RawURLEncoding.
			EncodeToString(
				[]byte(payload),
			) +
			"." +
			base64.
				RawURLEncoding.
				EncodeToString(
					signature,
				),
		nil
}

func ValidateTurnstileClearance(
	token string,
	clientKey string,
) bool {

	secret :=
		clearanceSecret()

	if secret == "" ||
		token == "" {
		return false
	}

	parts :=
		strings.Split(
			token,
			".",
		)

	if len(parts) != 2 {
		return false
	}

	payloadBytes, err :=
		base64.
			RawURLEncoding.
			DecodeString(
				parts[0],
			)

	if err != nil {
		return false
	}

	signature, err :=
		base64.
			RawURLEncoding.
			DecodeString(
				parts[1],
			)

	if err != nil {
		return false
	}

	payload :=
		string(
			payloadBytes,
		)

	expectedSignature :=
		signClearance(
			payload,
			secret,
		)

	if !hmac.Equal(
		signature,
		expectedSignature,
	) {
		return false
	}

	payloadParts :=
		strings.Split(
			payload,
			".",
		)

	if len(payloadParts) != 2 {
		return false
	}

	if payloadParts[0] !=
		hashClientKey(
			clientKey,
		) {

		return false
	}

	expiresAt, err :=
		strconv.ParseInt(
			payloadParts[1],
			10,
			64,
		)

	if err != nil {
		return false
	}

	return time.Now().Unix() <
		expiresAt
}

func TurnstileClearanceTTLSeconds() int {

	return configInt(
		"TurnstileClearanceTTLSeconds",
		600,
	)
}

func TurnstileCookieName() string {

	name :=
		strings.TrimSpace(
			web.AppConfig.
				DefaultString(
					"TurnstileCookieName",
					"atlas_turnstile_clearance",
				),
		)

	if name == "" ||
		strings.HasPrefix(
			name,
			"${",
		) {

		return "atlas_turnstile_clearance"
	}

	return name
}

func TurnstileCookieSecure() bool {

	value :=
		strings.ToLower(
			strings.TrimSpace(
				web.AppConfig.
					DefaultString(
						"TurnstileCookieSecure",
						"false",
					),
			),
		)

	return value == "true" ||
		value == "1" ||
		value == "yes"
}

func IsTurnstileDegraded() bool {

	turnstileAvailability.
		RLock()

	defer turnstileAvailability.
		RUnlock()

	return time.Now().
		Before(
			turnstileAvailability.
				DegradedUntil,
		)
}

func markTurnstileUnavailable() {

	degradedSeconds :=
		configInt(
			"TurnstileDegradedSeconds",
			60,
		)

	turnstileAvailability.
		Lock()

	turnstileAvailability.
		DegradedUntil =
		time.Now().
			Add(
				time.Duration(
					degradedSeconds,
				) * time.Second,
			)

	turnstileAvailability.
		Unlock()
}

func markTurnstileAvailable() {

	turnstileAvailability.
		Lock()

	turnstileAvailability.
		DegradedUntil =
		time.Time{}

	turnstileAvailability.
		Unlock()
}

func clearanceSecret() string {
	return getTurnstileConfig().
		ClearanceSecret
}

func hashClientKey(
	clientKey string,
) string {

	sum :=
		sha256.Sum256(
			[]byte(clientKey),
		)

	return hex.EncodeToString(
		sum[:],
	)
}

func signClearance(
	payload string,
	secret string,
) []byte {

	mac :=
		hmac.New(
			sha256.New,
			[]byte(secret),
		)

	_, _ =
		mac.Write(
			[]byte(payload),
		)

	return mac.Sum(nil)
}
