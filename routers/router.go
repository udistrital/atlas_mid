// @APIVersion 1.0.0
// @Title Atlas Externo Servicio API
// @Description API de consulta para Atlas Externo

package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/atlas_mid/controllers"
)

func init() {
	ns := web.NewNamespace(
		"/v1",

		web.NSNamespace(
			"/procesos",

			web.NSInclude(
				&controllers.ProcesoController{},
			),
		),

		web.NSNamespace(
			"/factores",
			web.NSInclude(
				&controllers.FactorController{},
			),
		),

		web.NSNamespace(
			"/caracteristicas",
			web.NSInclude(
				&controllers.CaracteristicaController{},
			),
		),

		web.NSNamespace(
			"/aspectos",

			web.NSInclude(
				&controllers.AspectoController{},
			),
		),

		web.NSNamespace(
			"/estructuras",

			web.NSInclude(
				&controllers.EstructuraController{},
			),
		),

		web.NSNamespace(
			"/security",

			web.NSInclude(
				&controllers.SecurityController{},
			),
		),
	)

	web.AddNamespace(ns)
}
