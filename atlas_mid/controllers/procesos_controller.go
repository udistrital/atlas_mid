package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/server/web"

	"github.com/udistrital/atlas_externo_servicio/services"
)

// ProcesoController operations for Proceso
type ProcesoController struct {
	web.Controller
}

// URLMapping ...
func (c *ProcesoController) URLMapping() {
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
// @Description consultar procesos
// @Success 200 {array} models.Proceso
// @Failure 500 error consultando procesos
// @router / [get]
func (c *ProcesoController) GetAll() {
	procesos, err := services.GetAllProcesos()

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		}

		c.Ctx.Output.SetStatus(
			http.StatusInternalServerError,
		)

		c.ServeJSON()

		return
	}

	c.Data["json"] = procesos

	c.ServeJSON()
}

// Get ...
// @Title Get
// @Description consultar proceso por id
// @Param id path string true "ID del proceso"
// @Success 200 {object} models.Proceso
// @Failure 404 proceso no encontrado
// @Failure 500 error consultando proceso
// @router /:id [get]
func (c *ProcesoController) Get() {
	id := c.Ctx.Input.Param(":id")

	proceso, err := services.GetProceso(id)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"err":    err.Error(),
			"status": http.StatusInternalServerError,
		}

		c.Ctx.Output.SetStatus(
			http.StatusInternalServerError,
		)

		c.ServeJSON()

		return
	}

	if proceso == nil {
		c.Data["json"] = map[string]interface{}{
			"err":    "proceso no encontrado",
			"status": http.StatusNotFound,
		}

		c.Ctx.Output.SetStatus(
			http.StatusNotFound,
		)

		c.ServeJSON()

		return
	}

	c.Data["json"] = proceso

	c.ServeJSON()
}
