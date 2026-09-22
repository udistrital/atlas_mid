package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/udistrital/atlas_externo_servicio/database"
	"github.com/udistrital/atlas_externo_servicio/models"
)

const procesoIndex = "atlas_procesos"

type elasticsearchProcesoHit struct {
	ID     string         `json:"_id"`
	Source models.Proceso `json:"_source"`
}

type elasticsearchProcesoSearchResponse struct {
	Hits struct {
		Hits []elasticsearchProcesoHit `json:"hits"`
	} `json:"hits"`
}

type elasticsearchProcesoGetResponse struct {
	ID     string         `json:"_id"`
	Source models.Proceso `json:"_source"`
}

func GetAllProcesos() (
	[]models.Proceso,
	error,
) {
	client, err := database.GetElasticsearchClient()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	query := `
	{
		"query": {
			"match_all": {}
		},
		"size": 1000
	}`

	response, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(
			procesoIndex,
		),
		client.Search.WithBody(
			strings.NewReader(query),
		),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando procesos: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando índice %s: %s",
			procesoIndex,
			response.Status(),
		)
	}

	var result elasticsearchProcesoSearchResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta Elasticsearch: %w",
			err,
		)
	}

	procesos := make(
		[]models.Proceso,
		0,
		len(result.Hits.Hits),
	)

	for _, hit := range result.Hits.Hits {
		proceso := hit.Source

		if proceso.ID == "" {
			proceso.ID = hit.ID
		}

		procesos = append(
			procesos,
			proceso,
		)
	}

	return procesos, nil
}

func GetProceso(
	id string,
) (*models.Proceso, error) {
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
		procesoIndex,
		id,
		client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error consultando proceso: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil, nil
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"error consultando proceso: %s",
			response.Status(),
		)
	}

	var result elasticsearchProcesoGetResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"error procesando respuesta Elasticsearch: %w",
			err,
		)
	}

	proceso := result.Source

	if proceso.ID == "" {
		proceso.ID = result.ID
	}

	return &proceso, nil
}
