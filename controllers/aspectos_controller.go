package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/atlas_mid/services"
)

type AspectoController struct {
	web.Controller
}

func (c *AspectoController) URLMapping() {
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
// @Description consultar aspectos
// @Param caracteristica_id query string false "ID de la característica"
// @Success 200 {array} models.Aspecto
// @Failure 500 error consultando aspectos
// @router / [get]
func (c *AspectoController) GetAll() {

	caracteristicaID := c.Ctx.Input.Query(
		"caracteristica_id",
	)

	aspectos, err :=
		services.GetAllAspectos(
			caracteristicaID,
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

	c.Data["json"] = aspectos

	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description consultar aspecto por id
// @Param id path string true "ID del aspecto"
// @Success 200 {object} models.Aspecto
// @Failure 404 aspecto no encontrado
// @Failure 500 error consultando aspecto
// @router /:id [get]
func (c *AspectoController) Get() {

	id := c.Ctx.Input.Param(":id")

	aspecto, err :=
		services.GetAspecto(id)

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

	if aspecto == nil {
		c.Data["json"] = map[string]interface{}{
			"error": "aspecto no encontrado",
		}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()
		return
	}

	c.Data["json"] = aspecto

	c.ServeJSON()
}
