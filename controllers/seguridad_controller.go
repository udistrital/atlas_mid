package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/atlas_mid/middleware"
	"github.com/udistrital/atlas_mid/services"
)

type SecurityController struct {
	web.Controller
}

type turnstileVerificationRequest struct {
	Token string `json:"token"`
}

func (
	c *SecurityController,
) URLMapping() {

	c.Mapping(
		"VerifyTurnstile",
		c.VerifyTurnstile,
	)
}

// VerifyTurnstile ...
// @Title VerifyTurnstile
// @Description valida un token Cloudflare Turnstile
// @Param body body turnstileVerificationRequest true "Token Turnstile"
// @Success 200 {object} map[string]interface{}
// @Failure 400 token inválido
// @Failure 429 demasiados intentos
// @Failure 503 Turnstile no disponible
// @router /turnstile/verify [post]
func (
	c *SecurityController,
) VerifyTurnstile() {

	var request turnstileVerificationRequest

	if err :=
		json.Unmarshal(
			c.Ctx.Input.RequestBody,
			&request,
		); err != nil {

		c.Data["json"] =
			map[string]interface{}{
				"success": false,
				"message": "Solicitud inválida.",
			}

		c.Ctx.Output.
			SetStatus(
				http.StatusBadRequest,
			)

		c.ServeJSON()

		return
	}

	clientKey :=
		middleware.ClientKey(
			c.Ctx,
		)

	if !services.
		AllowTurnstileVerification(
			clientKey,
		) {

		c.Data["json"] =
			map[string]interface{}{
				"success": false,
				"code":    "TURNSTILE_RATE_LIMITED",
				"message": "Demasiados intentos de verificación.",
			}

		c.Ctx.Output.
			SetStatus(
				http.StatusTooManyRequests,
			)

		c.ServeJSON()

		return
	}

	result, err :=
		services.ValidateTurnstile(
			request.Token,
			extractIP(
				clientKey,
			),
		)

	if err != nil {

		if errors.Is(
			err,
			services.
				ErrTurnstileUnavailable,
		) {

			c.Data["json"] =
				map[string]interface{}{
					"success": false,
					"code":    "TURNSTILE_UNAVAILABLE",
					"message": "La verificación de seguridad no está disponible temporalmente.",
				}

			c.Ctx.Output.
				SetStatus(
					http.StatusServiceUnavailable,
				)

			c.ServeJSON()

			return
		}

		c.Data["json"] =
			map[string]interface{}{
				"success": false,
				"message": "No fue posible validar la verificación.",
			}

		c.Ctx.Output.
			SetStatus(
				http.StatusInternalServerError,
			)

		c.ServeJSON()

		return
	}

	if result == nil ||
		!result.Success {

		if result != nil {

			logs.Warn(
				"Turnstile inválido: hostname=%q action=%q errores=%v",
				result.Hostname,
				result.Action,
				result.ErrorCodes,
			)
		}

		c.Data["json"] =
			map[string]interface{}{
				"success": false,
				"code":    "TURNSTILE_INVALID",
				"message": "La verificación de seguridad no fue válida.",
			}

		c.Ctx.Output.
			SetStatus(
				http.StatusBadRequest,
			)

		c.ServeJSON()

		return
	}

	clearance, err :=
		services.
			GenerateTurnstileClearance(
				clientKey,
			)

	if err != nil {

		c.Data["json"] =
			map[string]interface{}{
				"success": false,
				"message": "No fue posible generar la autorización temporal.",
			}

		c.Ctx.Output.
			SetStatus(
				http.StatusInternalServerError,
			)

		c.ServeJSON()

		return
	}

	c.Ctx.SetCookie(
		services.
			TurnstileCookieName(),
		clearance,
		services.
			TurnstileClearanceTTLSeconds(),
		"/",
		"",
		services.
			TurnstileCookieSecure(),
		true,
	)

	services.ResetRisk(
		clientKey,
	)

	c.Data["json"] =
		map[string]interface{}{
			"success": true,

			"message": "Verificación completada correctamente.",
		}

	c.ServeJSON()
}

func extractIP(
	clientKey string,
) string {

	for i, value := range clientKey {

		if value == '|' {
			return clientKey[:i]
		}
	}

	return clientKey
}
