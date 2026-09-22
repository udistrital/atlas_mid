package models

type Factor struct {
	ID                string   `json:"id"`
	ProcesoID         string   `json:"proceso_id"`
	Nombre            string   `json:"nombre"`
	Descripcion       string   `json:"descripcion"`
	Calificacion      *float64 `json:"calificacion"`
	Caracteristicas   []string `json:"caracteristicas"`
	Activo            bool     `json:"activo"`
	FechaCreacion     *string  `json:"fecha_creacion"`
	FechaModificacion *string  `json:"fecha_modificacion"`
}
