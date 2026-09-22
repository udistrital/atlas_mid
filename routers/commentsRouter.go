package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:AspectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:AspectoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:AspectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:AspectoController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:CaracteristicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:CaracteristicaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:CaracteristicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:CaracteristicaController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:EstructuraController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:EstructuraController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:EstructuraController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:EstructuraController"],
        beego.ControllerComments{
            Method: "GetDatos",
            Router: `/:id/datos`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:FactorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:FactorController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:FactorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:FactorController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:ProcesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:ProcesoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:ProcesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:ProcesoController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:SecurityController"] = append(beego.GlobalControllerRouter["github.com/udistrital/atlas_externo_servicio/controllers:SecurityController"],
        beego.ControllerComments{
            Method: "VerifyTurnstile",
            Router: `/turnstile/verify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
