package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchClient(host string, port string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", host, port),
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	return &ElasticsearchClient{
		client: client,
	}, nil
}

// IndexDocument indexes a document in Elasticsearch
func (es *ElasticsearchClient) IndexDocument(ctx context.Context, index string, document interface{}) error {
	docBytes, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	res, err := es.client.Index(index, bytes.NewReader(docBytes))
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// SearchDocuments performs a search query in Elasticsearch
func (es *ElasticsearchClient) SearchDocuments(ctx context.Context, index string, query map[string]interface{}) ([]json.RawMessage, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(index),
		es.client.Search.WithBody(bytes.NewReader(queryBytes)),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]json.RawMessage, len(hits))
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		docBytes, err := json.Marshal(source)
		if err != nil {
			return nil, fmt.Errorf("error marshaling document: %w", err)
		}
		documents[i] = docBytes
	}

	return documents, nil
}
