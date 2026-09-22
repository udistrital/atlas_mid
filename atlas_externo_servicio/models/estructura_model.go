package models

type Estructura struct {
	ID            string                   `json:"id"`
	AspectoID     string                   `json:"aspecto_id"`
	TipoEvidencia string                   `json:"tipo_evidencia"`
	Nombre        string                   `json:"nombre"`
	Activo        bool                     `json:"activo"`
	Campos        []map[string]interface{} `json:"campos"`
	Data          []map[string]interface{} `json:"data"`
}
