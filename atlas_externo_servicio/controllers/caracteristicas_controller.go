package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/atlas_externo_servicio/services"
)

type CaracteristicaController struct {
	web.Controller
}

func (c *CaracteristicaController) URLMapping() {
	c.Mapping(
		"GetAll",
		c.GetAll,
	)

	c.Mapping(
		"Get",
		c.Get,
	)
}

// GetAll ...
// @Title GetAll
// @Description consultar características
// @Param factor_id query string false "ID del factor"
// @Success 200 {array} models.Caracteristica
// @Failure 500 error consultando características
// @router / [get]
func (c *CaracteristicaController) GetAll() {

	factorID := c.Ctx.Input.Query("factor_id")

	caracteristicas, err :=
		services.GetAllCaracteristicas(factorID)

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

	c.Data["json"] = caracteristicas
	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description consultar característica por id
// @Param id path string true "ID de la característica"
// @Success 200 {object} models.Caracteristica
// @Failure 404 característica no encontrada
// @Failure 500 error consultando característica
// @router /:id [get]
func (c *CaracteristicaController) Get() {

	id := c.Ctx.Input.Param(":id")

	caracteristica, err :=
		services.GetCaracteristica(id)

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

	if caracteristica == nil {
		c.Data["json"] = map[string]interface{}{
			"error": "característica no encontrada",
		}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()
		return
	}

	c.Data["json"] = caracteristica
	c.ServeJSON()
}
