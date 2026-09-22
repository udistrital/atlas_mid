package models

type Proceso struct {
	ID                     string   `json:"id"`
	Nombre                 string   `json:"nombre"`
	Descripcion            string   `json:"descripcion"`
	DependenciaResponsable string   `json:"dependencia_responsable"`
	Objetivo               string   `json:"objetivo"`
	Factores               []string `json:"factores"`
	FechaInicio            *string  `json:"fecha_inicio"`
	FechaFin               *string  `json:"fecha_fin"`
	Activo                 bool     `json:"activo"`
	FechaCreacion          *string  `json:"fecha_creacion"`
	FechaModificacion      *string  `json:"fecha_modificacion"`
}
