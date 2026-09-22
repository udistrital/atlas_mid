package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/udistrital/atlas_mid/database"
	"github.com/udistrital/atlas_mid/models"
)

type elasticsearchEstructuraGetResponse struct {
	ID     string            `json:"_id"`
	Source models.Estructura `json:"_source"`
}

type DatosEstructuraResponse struct {
	Count    int                      `json:"count"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Results  []map[string]interface{} `json:"results"`
}

func GetEstructura(
	id string,
) (*models.Estructura, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, fmt.Errorf(
			"el id de la estructura es obligatorio",
		)
	}

	if !esIndiceEstructuraValido(id) {
		return nil, fmt.Errorf(
			"índice de estructura no válido: %s",
			id,
		)
	}

	client, err :=
		database.GetElasticsearchClient()

	if err != nil {
		return nil, err
	}

	ctx, cancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

	defer cancel()

	/*
		En el CRUD Django el campo id de la
		estructura contiene el nombre completo
		del índice Elasticsearch.

		Ejemplo:

		atlas_estructura_documental_UUID
	*/
	indexName := id

	response, err := client.Get(
		indexName,
		id,
		client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando estructura %s: %w",
			id,
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode ==
		http.StatusNotFound {

		return nil, nil
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando estructura %s: %s",
			id,
			response.Status(),
		)
	}

	var result elasticsearchEstructuraGetResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	estructura := result.Source

	if estructura.ID == "" {
		estructura.ID = result.ID
	}

	if estructura.Campos == nil {
		estructura.Campos =
			[]map[string]interface{}{}
	}

	if estructura.Data == nil {
		estructura.Data =
			[]map[string]interface{}{}
	}

	return &estructura, nil
}

func GetDatosEstructura(
	id string,
	query url.Values,
) (*DatosEstructuraResponse, error) {

	estructura, err :=
		GetEstructura(id)

	if err != nil {
		return nil, err
	}

	if estructura == nil {
		return nil, nil
	}

	camposActivos :=
		obtenerCamposActivos(
			estructura.Campos,
		)

	filas := make(
		[]map[string]interface{},
		0,
		len(estructura.Data),
	)

	for index, item := range estructura.Data {

		if item == nil {
			continue
		}

		fila :=
			map[string]interface{}{
				"id": strconv.Itoa(index),
				"activo": obtenerValor(
					item,
					"activo",
					true,
				),
				"fecha_creacion":     item["fecha_creacion"],
				"fecha_modificacion": item["fecha_modificacion"],
			}

		valores :=
			obtenerMapa(
				item["valores"],
			)

		for _, campo := range camposActivos {

			campoID :=
				obtenerString(
					campo["campo_id"],
				)

			nombreCampo :=
				obtenerString(
					campo["nombre_campo"],
				)

			if campoID == "" ||
				nombreCampo == "" {
				continue
			}

			if valor, existe :=
				valores[campoID]; existe {

				fila[nombreCampo] =
					valor

				continue
			}

			/*
				Compatibilidad con datos
				antiguos guardados directamente
				con nombre_campo.
			*/
			if valor, existe :=
				item[nombreCampo]; existe {

				fila[nombreCampo] =
					valor

				continue
			}

			fila[nombreCampo] = nil
		}

		filas = append(
			filas,
			fila,
		)
	}

	filas =
		aplicarFiltrosDatos(
			filas,
			query,
		)

	aplicarOrdenamientoDatos(
		filas,
		query.Get("ordering"),
	)

	page :=
		obtenerEnteroPositivo(
			query.Get("page"),
			1,
		)

	pageSize :=
		obtenerEnteroPositivo(
			query.Get("page_size"),
			10,
		)

	total := len(filas)

	inicio :=
		(page - 1) *
			pageSize

	if inicio > total {
		inicio = total
	}

	fin :=
		inicio +
			pageSize

	if fin > total {
		fin = total
	}

	results :=
		filas[inicio:fin]

	if results == nil {
		results =
			[]map[string]interface{}{}
	}

	return &DatosEstructuraResponse{
		Count:    total,
		Page:     page,
		PageSize: pageSize,
		Results:  results,
	}, nil
}

func esIndiceEstructuraValido(
	id string,
) bool {

	return strings.HasPrefix(
		id,
		"atlas_estructura_documental_",
	) ||
		strings.HasPrefix(
			id,
			"atlas_estructura_tabla_",
		)
}

func obtenerCamposActivos(
	campos []map[string]interface{},
) []map[string]interface{} {

	resultado :=
		make(
			[]map[string]interface{},
			0,
			len(campos),
		)

	for _, campo := range campos {

		if campo == nil {
			continue
		}

		if activo,
			existe :=
			campo["activo"]; existe {

			if valor,
				ok :=
				activo.(bool); ok &&
				!valor {

				continue
			}
		}

		resultado = append(
			resultado,
			campo,
		)
	}

	sort.SliceStable(
		resultado,
		func(i, j int) bool {

			return obtenerOrden(
				resultado[i],
			) <
				obtenerOrden(
					resultado[j],
				)
		},
	)

	return resultado
}

func obtenerOrden(
	campo map[string]interface{},
) float64 {

	switch valor :=
		campo["orden"].(type) {

	case float64:
		if valor > 0 {
			return valor
		}

	case int:
		if valor > 0 {
			return float64(valor)
		}
	}

	return 999999
}

func obtenerMapa(
	value interface{},
) map[string]interface{} {

	if value == nil {
		return map[string]interface{}{}
	}

	resultado, ok := value.(map[string]interface{})

	if !ok {
		return map[string]interface{}{}
	}

	return resultado
}

func obtenerString(
	value interface{},
) string {

	if value == nil {
		return ""
	}

	return strings.TrimSpace(
		fmt.Sprint(value),
	)
}

func obtenerValor(
	data map[string]interface{},
	key string,
	defaultValue interface{},
) interface{} {

	value, existe :=
		data[key]

	if !existe ||
		value == nil {

		return defaultValue
	}

	return value
}

func obtenerEnteroPositivo(
	value string,
	defaultValue int,
) int {

	numero, err :=
		strconv.Atoi(value)

	if err != nil ||
		numero < 1 {

		return defaultValue
	}

	return numero
}

func aplicarFiltrosDatos(
	filas []map[string]interface{},
	query url.Values,
) []map[string]interface{} {

	excluidos :=
		map[string]bool{
			"page":      true,
			"page_size": true,
			"ordering":  true,
			"format":    true,
		}

	resultado := filas

	for clave, valores := range query {

		if excluidos[clave] ||
			len(valores) == 0 {

			continue
		}

		busqueda :=
			strings.ToLower(
				strings.TrimSpace(
					valores[0],
				),
			)

		if busqueda == "" {
			continue
		}

		filtradas :=
			make(
				[]map[string]interface{},
				0,
			)

		for _, fila := range resultado {

			valor :=
				strings.ToLower(
					fmt.Sprint(
						fila[clave],
					),
				)

			if strings.Contains(
				valor,
				busqueda,
			) {
				filtradas = append(
					filtradas,
					fila,
				)
			}
		}

		resultado = filtradas
	}

	return resultado
}

func aplicarOrdenamientoDatos(
	filas []map[string]interface{},
	ordering string,
) {

	ordering =
		strings.TrimSpace(
			ordering,
		)

	if ordering == "" {
		return
	}

	descendente :=
		strings.HasPrefix(
			ordering,
			"-",
		)

	campo :=
		strings.TrimPrefix(
			ordering,
			"-",
		)

	if campo == "" {
		return
	}

	sort.SliceStable(
		filas,
		func(i, j int) bool {

			valorI :=
				filas[i][campo]

			valorJ :=
				filas[j][campo]

			if valorI == nil &&
				valorJ == nil {
				return false
			}

			if valorI == nil {
				return descendente
			}

			if valorJ == nil {
				return !descendente
			}

			textoI :=
				strings.ToLower(
					fmt.Sprint(
						valorI,
					),
				)

			textoJ :=
				strings.ToLower(
					fmt.Sprint(
						valorJ,
					),
				)

			if descendente {
				return textoI >
					textoJ
			}

			return textoI <
				textoJ
		},
	)
}
