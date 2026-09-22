package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/udistrital/atlas_mid/database"
	"github.com/udistrital/atlas_mid/models"
)

const caracteristicaIndex = "atlas_caracteristicas"

type elasticsearchCaracteristicaHit struct {
	ID     string                `json:"_id"`
	Source models.Caracteristica `json:"_source"`
}

type elasticsearchCaracteristicaSearchResponse struct {
	Hits struct {
		Hits []elasticsearchCaracteristicaHit `json:"hits"`
	} `json:"hits"`
}

type elasticsearchCaracteristicaGetResponse struct {
	ID     string                `json:"_id"`
	Source models.Caracteristica `json:"_source"`
}

func GetAllCaracteristicas(
	factorID string,
) ([]models.Caracteristica, error) {

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

	if factorID != "" {
		query = map[string]interface{}{
			"term": map[string]interface{}{
				"factor_id": factorID,
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
			"error construyendo consulta de características: %w",
			err,
		)
	}

	response, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(caracteristicaIndex),
		client.Search.WithBody(&body),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando características: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando índice %s: %s",
			caracteristicaIndex,
			response.Status(),
		)
	}

	var result elasticsearchCaracteristicaSearchResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	caracteristicas := make(
		[]models.Caracteristica,
		0,
		len(result.Hits.Hits),
	)

	for _, hit := range result.Hits.Hits {
		caracteristica := hit.Source

		if caracteristica.ID == "" {
			caracteristica.ID = hit.ID
		}

		if caracteristica.Aspectos == nil {
			caracteristica.Aspectos = []string{}
		}

		caracteristicas = append(
			caracteristicas,
			caracteristica,
		)
	}

	return caracteristicas, nil
}

func GetCaracteristica(
	id string,
) (*models.Caracteristica, error) {

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
		caracteristicaIndex,
		id,
		client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando característica: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando característica: %s",
			response.Status(),
		)
	}

	var result elasticsearchCaracteristicaGetResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	caracteristica := result.Source

	if caracteristica.ID == "" {
		caracteristica.ID = result.ID
	}

	if caracteristica.Aspectos == nil {
		caracteristica.Aspectos = []string{}
	}

	return &caracteristica, nil
}
