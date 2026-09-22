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

const factorIndex = "atlas_factores"

type elasticsearchFactorHit struct {
	ID     string        `json:"_id"`
	Source models.Factor `json:"_source"`
}

type elasticsearchFactorSearchResponse struct {
	Hits struct {
		Hits []elasticsearchFactorHit `json:"hits"`
	} `json:"hits"`
}

type elasticsearchFactorGetResponse struct {
	ID     string        `json:"_id"`
	Source models.Factor `json:"_source"`
}

func GetAllFactores(
	procesoID string,
) ([]models.Factor, error) {

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

	if procesoID != "" {
		query = map[string]interface{}{
			"term": map[string]interface{}{
				"proceso_id": procesoID,
			},
		}
	}

	requestBody := map[string]interface{}{
		"query": query,
		"size":  1000,
	}

	var body bytes.Buffer

	if err := json.NewEncoder(
		&body,
	).Encode(requestBody); err != nil {

		return nil, fmt.Errorf(
			"error construyendo consulta de factores: %w",
			err,
		)
	}

	response, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(
			factorIndex,
		),
		client.Search.WithBody(
			&body,
		),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando factores: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando índice %s: %s",
			factorIndex,
			response.Status(),
		)
	}

	var result elasticsearchFactorSearchResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	factores := make(
		[]models.Factor,
		0,
		len(result.Hits.Hits),
	)

	for _, hit := range result.Hits.Hits {

		factor := hit.Source

		if factor.ID == "" {
			factor.ID = hit.ID
		}

		if factor.Caracteristicas == nil {
			factor.Caracteristicas = []string{}
		}

		factores = append(
			factores,
			factor,
		)
	}

	return factores, nil
}

func GetFactor(
	id string,
) (*models.Factor, error) {

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
		factorIndex,
		id,
		client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando factor: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando factor: %s",
			response.Status(),
		)
	}

	var result elasticsearchFactorGetResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return nil, fmt.Errorf(
			"error procesando respuesta de Elasticsearch: %w",
			err,
		)
	}

	factor := result.Source

	if factor.ID == "" {
		factor.ID = result.ID
	}

	if factor.Caracteristicas == nil {
		factor.Caracteristicas = []string{}
	}

	return &factor, nil
}
