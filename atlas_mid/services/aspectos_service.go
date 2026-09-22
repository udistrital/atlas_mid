package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/udistrital/atlas_externo_servicio/database"
	"github.com/udistrital/atlas_externo_servicio/models"
)

const aspectoIndex = "atlas_aspectos"

type elasticsearchAspectoHit struct {
	ID     string         `json:"_id"`
	Source models.Aspecto `json:"_source"`
}

type elasticsearchAspectoSearchResponse struct {
	Hits struct {
		Hits []elasticsearchAspectoHit `json:"hits"`
	} `json:"hits"`
}

type elasticsearchAspectoGetResponse struct {
	ID     string         `json:"_id"`
	Source models.Aspecto `json:"_source"`
}

func GetAllAspectos(
	caracteristicaID string,
) ([]models.Aspecto, error) {

	client, err := database.GetElasticsearchClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	query := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}

	if caracteristicaID != "" {
		query = map[string]interface{}{
			"term": map[string]interface{}{
				"caracteristica_id": caracteristicaID,
			},
		}
	}

	requestBody := map[string]interface{}{
		"query": query,
		"size":  1000,
	}

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(requestBody); err != nil {
		return nil, fmt.Errorf(
			"error construyendo consulta de aspectos: %w",
			err,
		)
	}

	response, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(aspectoIndex),
		client.Search.WithBody(&body),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando aspectos: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando índice %s: %s",
			aspectoIndex,
			response.Status(),
		)
	}

	var result elasticsearchAspectoSearchResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	aspectos := make(
		[]models.Aspecto,
		0,
		len(result.Hits.Hits),
	)

	for _, hit := range result.Hits.Hits {
		aspecto := hit.Source

		if aspecto.ID == "" {
			aspecto.ID = hit.ID
		}

		if aspecto.EstructurasEvidencias == nil {
			aspecto.EstructurasEvidencias =
				[]map[string]interface{}{}
		}

		aspectos = append(
			aspectos,
			aspecto,
		)
	}

	return aspectos, nil
}

func GetAspecto(
	id string,
) (*models.Aspecto, error) {

	client, err := database.GetElasticsearchClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	response, err := client.Get(
		aspectoIndex,
		id,
		client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando aspecto: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando aspecto: %s",
			response.Status(),
		)
	}

	var result elasticsearchAspectoGetResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	aspecto := result.Source

	if aspecto.ID == "" {
		aspecto.ID = result.ID
	}

	if aspecto.EstructurasEvidencias == nil {
		aspecto.EstructurasEvidencias =
			[]map[string]interface{}{}
	}

	return &aspecto, nil
}
