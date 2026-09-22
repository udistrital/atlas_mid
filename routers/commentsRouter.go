package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:AspectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:AspectoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:AspectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:AspectoController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:CaracteristicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:CaracteristicaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:CaracteristicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:CaracteristicaController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:EstructuraController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:EstructuraController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:EstructuraController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:EstructuraController"],
		beego.ControllerComments{
			Method:           "GetDatos",
			Router:           `/:id/datos`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:FactorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:FactorController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:FactorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:FactorController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:ProcesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:ProcesoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:ProcesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:ProcesoController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:SecurityController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_mid/controllers:SecurityController"],
		beego.ControllerComments{
			Method:           "VerifyTurnstile",
			Router:           `/turnstile/verify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
