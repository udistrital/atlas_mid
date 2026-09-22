package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/udistrital/atlas_mid/services"
)

func RiskMiddleware(
	ctx *beecontext.Context,
) {

	path :=
		ctx.Request.URL.Path

	/*
		El endpoint de Turnstile debe poder
		ejecutarse incluso cuando el cliente
		está siendo desafiado.
	*/
	if strings.HasPrefix(
		path,
		"/v1/security/",
	) {
		return
	}

	clientKey :=
		ClientKey(ctx)

	decision :=
		services.EvaluateRisk(
			clientKey,
			path,
		)

	logs.Info(
		"risk path=%s score=%d requests=%d burst=%d unique=%d challenge=%v block=%v",
		path,
		decision.Score,
		decision.RequestsInWindow,
		decision.RequestsInBurst,
		decision.UniquePathsWindow,
		decision.RequireTurnstile,
		decision.Block,
	)

	if decision.Block {

		writeJSON(
			ctx,
			http.StatusTooManyRequests,
			map[string]interface{}{
				"success": false,
				"code":    "RATE_LIMITED",
				"message": "Se excedió temporalmente el límite de consultas.",
			},
		)

		return
	}

	clearance :=
		ctx.GetCookie(
			services.
				TurnstileCookieName(),
		)

	hasValidClearance :=
		services.
			ValidateTurnstileClearance(
				clearance,
				clientKey,
			)

	/*
		Si ya superó Turnstile,
		normalmente dejamos continuar.

		Pero si vuelve a alcanzar un
		nivel alto de riesgo, exigimos
		un nuevo desafío incluso antes
		de que expire la clearance.
	*/
	if hasValidClearance {

		if decision.Score <
			services.
				RiskRechallengeScore() {

			return
		}

		/*
			Si Cloudflare está caído,
			no convertimos Turnstile en
			un punto único de fallo.

			El bloqueo crítico ya fue
			evaluado arriba.
		*/
		if services.
			IsTurnstileDegraded() {

			return
		}

		/*
			Invalidamos la clearance
			actual para obligar a generar
			una nueva.
		*/
		ctx.SetCookie(
			services.
				TurnstileCookieName(),
			"",
			-1,
			"/",
			"",
			services.
				TurnstileCookieSecure(),
			true,
		)

		writeJSON(
			ctx,
			http.StatusPreconditionRequired,
			map[string]interface{}{
				"success": false,

				"code": "TURNSTILE_REQUIRED",

				"message": "Se requiere una nueva verificación de seguridad.",
			},
		)

		return
	}

	/*
		Sin clearance:
		solo pedimos Turnstile si
		el score normal alcanzó el
		umbral de challenge.
	*/
	if !decision.RequireTurnstile {
		return
	}

	if services.
		IsTurnstileDegraded() {

		return
	}

	writeJSON(
		ctx,
		http.StatusPreconditionRequired,
		map[string]interface{}{
			"success": false,

			"code": "TURNSTILE_REQUIRED",

			"message": "Se requiere una verificación de seguridad para continuar.",
		},
	)
}

func ClientKey(
	ctx *beecontext.Context,
) string {

	ip :=
		clientIP(ctx)

	/*
		WSO2 ya valida el Bearer.

		No guardamos el Bearer.
		Solo usamos su hash como señal adicional
		para diferenciar clientes que puedan
		compartir una IP/NAT.
	*/
	authorization :=
		ctx.Request.
			Header.
			Get(
				"Authorization",
			)

	if authorization == "" {
		return ip
	}

	sum :=
		sha256.Sum256(
			[]byte(
				authorization,
			),
		)

	tokenHash :=
		hex.EncodeToString(
			sum[:],
		)

	return ip +
		"|" +
		tokenHash
}

func clientIP(
	ctx *beecontext.Context,
) string {

	forwardedFor :=
		ctx.Request.
			Header.
			Get(
				"X-Forwarded-For",
			)

	if forwardedFor != "" {

		first :=
			strings.TrimSpace(
				strings.Split(
					forwardedFor,
					",",
				)[0],
			)

		if net.ParseIP(
			first,
		) != nil {

			return first
		}
	}

	host, _, err :=
		net.SplitHostPort(
			ctx.Request.
				RemoteAddr,
		)

	if err == nil &&
		host != "" {

		return host
	}

	if ctx.Request.
		RemoteAddr != "" {

		return ctx.Request.
			RemoteAddr
	}

	return "unknown"
}

func writeJSON(
	ctx *beecontext.Context,
	status int,
	data interface{},
) {

	ctx.Output.
		SetStatus(
			status,
		)

	_ =
		ctx.JSONResp(
			data,
		)
}
