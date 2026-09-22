package models

type Caracteristica struct {
	ID                string   `json:"id"`
	FactorID          string   `json:"factor_id"`
	Nombre            string   `json:"nombre"`
	Descripcion       string   `json:"descripcion"`
	Calificacion      *float64 `json:"calificacion"`
	Aspectos          []string `json:"aspectos"`
	Activo            bool     `json:"activo"`
	FechaCreacion     *string  `json:"fecha_creacion"`
	FechaModificacion *string  `json:"fecha_modificacion"`
}
