package main

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"

	"github.com/udistrital/atlas_mid/database"
	_ "github.com/udistrital/atlas_mid/routers"

	"github.com/udistrital/atlas_mid/middleware"
	"github.com/udistrital/atlas_mid/services"
	apistatus "github.com/udistrital/utils_oas/v2/apiStatusLib"
	"github.com/udistrital/utils_oas/v2/auditoria"
	customerrorv2 "github.com/udistrital/utils_oas/v2/customerror"
	"github.com/udistrital/utils_oas/v2/security"
	"github.com/udistrital/utils_oas/v2/xray"
)

func init() {
	//Elasticsearch
	if err :=
		database.
			InitElasticsearch(); err != nil {

		logs.Critical(
			"error conectando a Elasticsearch: %v",
			err,
		)

		panic(err)
	}

	logs.Info(
		"conexión a Elasticsearch establecida",
	)
	// Turnstile
	if err :=
		services.
			ValidateTurnstileConfiguration(); err != nil {

		logs.Critical(
			"configuración Turnstile inválida: %v",
			err,
		)

		panic(err)
	}

	logs.Info(
		"configuración Turnstile validada",
	)

	// CORS
	allowedOrigins := []string{
		"*.udistrital.edu.co",
	}

	if web.BConfig.RunMode == web.DEV {
		allowedOrigins = []string{
			"*",
		}
	}

	web.InsertFilter(
		"*",
		web.BeforeRouter,
		cors.Allow(
			&cors.Options{
				AllowOrigins: allowedOrigins,

				AllowMethods: []string{
					"GET",
					"POST",
					"OPTIONS",
				},

				AllowHeaders: []string{
					"Accept",
					"Authorization",
					"Content-Type",
					"User-Agent",
					"X-Amzn-Trace-Id",
				},

				ExposeHeaders: []string{
					"Content-Length",
				},

				AllowCredentials: true,
			},
		),
	)

	web.InsertFilter(
		"/v1/*",
		web.BeforeExec,
		middleware.RiskMiddleware,
	)

	// Swagger
	if web.AppConfig.DefaultBool(
		"EnableDocs",
		false,
	) {
		web.SetStaticPath(
			"/swagger",
			"swagger",
		)
	}

	// Integraciones institucionales
	apistatus.Init()
	auditoria.InitMiddleware()
	security.SetSecurityHeaders()
	xray.Init()

	web.ErrorController(
		&customerrorv2.CustomErrorController{},
	)
}

func main() {
	web.Run()
}
