package models

type Aspecto struct {
	ID                    string                   `json:"id"`
	CaracteristicaID      string                   `json:"caracteristica_id"`
	Nombre                string                   `json:"nombre"`
	EstructurasEvidencias []map[string]interface{} `json:"estructuras_evidencias"`
	Activo                bool                     `json:"activo"`
	FechaCreacion         *string                  `json:"fecha_creacion"`
	FechaModificacion     *string                  `json:"fecha_modificacion"`
}
