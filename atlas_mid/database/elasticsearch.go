package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/elastic/go-elasticsearch/v8"
)

var elasticsearchClient *elasticsearch.Client

func InitElasticsearch() error {
	scheme := web.AppConfig.DefaultString(
		"ElasticsearchScheme",
		"http",
	)

	host := web.AppConfig.DefaultString(
		"ElasticsearchHost",
		"localhost",
	)

	port := web.AppConfig.DefaultString(
		"ElasticsearchPort",
		"9200",
	)

	username := web.AppConfig.DefaultString(
		"ElasticsearchUsername",
		"",
	)

	password := web.AppConfig.DefaultString(
		"ElasticsearchPassword",
		"",
	)

	address := fmt.Sprintf(
		"%s://%s:%s",
		scheme,
		host,
		port,
	)

	cfg := elasticsearch.Config{
		Addresses: []string{
			address,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			MaxIdleConns:          20,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf(
			"error creando cliente Elasticsearch: %w",
			err,
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	/*
		No usamos client.Info() porque el usuario
		app_atlas_externo no tiene permisos de cluster.

		Validamos la conexión contra un índice al cual
		sí tiene permisos de lectura/metadatos.
	*/
	response, err := client.Indices.Exists(
		[]string{"atlas_procesos"},
		client.Indices.Exists.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf(
			"error validando conexión con Elasticsearch: %w",
			err,
		)
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		elasticsearchClient = client
		return nil

	case http.StatusUnauthorized:
		return errors.New(
			"credenciales de Elasticsearch inválidas",
		)

	case http.StatusForbidden:
		return errors.New(
			"usuario sin permisos sobre el índice atlas_procesos",
		)

	case http.StatusNotFound:
		return errors.New(
			"índice atlas_procesos no encontrado",
		)

	default:
		return fmt.Errorf(
			"Elasticsearch respondió con error: %s",
			response.Status(),
		)
	}
}

func GetElasticsearchClient() (
	*elasticsearch.Client,
	error,
) {
	if elasticsearchClient == nil {
		return nil, errors.New(
			"cliente Elasticsearch no inicializado",
		)
	}

	return elasticsearchClient, nil
}
