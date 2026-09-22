package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/atlas_externo_servicio/services"
)

type EstructuraController struct {
	web.Controller
}

func (c *EstructuraController) URLMapping() {
	c.Mapping(
		"Get",
		c.Get,
	)

	c.Mapping(
		"GetDatos",
		c.GetDatos,
	)
}

// Get ...
// @Title Get
// @Description consultar estructura por id
// @Param id path string true "ID de la estructura"
// @Success 200 {object} models.Estructura
// @Failure 404 estructura no encontrada
// @Failure 500 error consultando estructura
// @router /:id [get]
func (c *EstructuraController) Get() {

	id := c.Ctx.Input.Param(":id")

	estructura, err :=
		services.GetEstructura(id)

	if err != nil {
		c.Data["json"] =
			map[string]interface{}{
				"error": err.Error(),
			}

		c.Ctx.Output.SetStatus(
			http.StatusInternalServerError,
		)

		c.ServeJSON()
		return
	}

	if estructura == nil {
		c.Data["json"] =
			map[string]interface{}{
				"error": "estructura no encontrada",
			}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()
		return
	}

	c.Data["json"] = estructura
	c.ServeJSON()
}

// GetDatos ...
// @Title GetDatos
// @Description consultar datos paginados de una estructura
// @Param id path string true "ID de la estructura"
// @Param page query int false "Página"
// @Param page_size query int false "Tamaño de página"
// @Param ordering query string false "Campo de ordenamiento"
// @Success 200 {object} services.DatosEstructuraResponse
// @Failure 404 estructura no encontrada
// @Failure 500 error consultando datos
// @router /:id/datos [get]
func (c *EstructuraController) GetDatos() {
	id := c.Ctx.Input.Param(":id")

	resultado, err := services.GetDatosEstructura(
		id,
		c.Ctx.Request.URL.Query(),
	)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": err.Error(),
		}

		c.Ctx.Output.SetStatus(
			http.StatusInternalServerError,
		)

		c.ServeJSON()
		return
	}

	if resultado == nil {
		c.Data["json"] = map[string]interface{}{
			"error": "estructura no encontrada",
		}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()
		return
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}
