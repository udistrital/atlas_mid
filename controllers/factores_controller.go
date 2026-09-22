package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"

	"github.com/udistrital/atlas_externo_servicio/services"
)

type FactorController struct {
	web.Controller
}

func (c *FactorController) URLMapping() {

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
// @Description consultar factores
// @Param proceso_id query string false "ID del proceso"
// @Success 200 {array} models.Factor
// @Failure 500 error consultando factores
// @router / [get]
func (c *FactorController) GetAll() {

	procesoID := c.Ctx.Input.Query(
		"proceso_id",
	)

	factores, err := services.GetAllFactores(
		procesoID,
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

	c.Data["json"] = factores

	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description consultar factor por id
// @Param id path string true "ID del factor"
// @Success 200 {object} models.Factor
// @Failure 404 factor no encontrado
// @Failure 500 error consultando factor
// @router /:id [get]
func (c *FactorController) Get() {

	id := c.Ctx.Input.Param(":id")

	factor, err := services.GetFactor(id)

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

	if factor == nil {

		c.Data["json"] = map[string]interface{}{
			"error": "factor no encontrado",
		}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()

		return
	}

	c.Data["json"] = factor

	c.ServeJSON()
}
