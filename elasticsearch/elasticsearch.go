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
func (es *ElasticsearchClient) SearchDocuments(ctx context.Context, query *ESQuery) ([]json.RawMessage, error) {
	if query == nil {
		return nil, fmt.Errorf("query is nil")
	}

	if query.query == nil {
		return nil, fmt.Errorf("query is nil")
	}

	queryBytes, err := json.Marshal(query.query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(query.index),
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

// SearchDocumentsWithQuery performs a multi-search query in Elasticsearch
func (es *ElasticsearchClient) SearchDocumentsWithQuery(ctx context.Context, index string, query *MultiESQuery) ([][]json.RawMessage, error) {
	buffer, err := query.createMQueryBuffer(index)
	if err != nil {
		return nil, fmt.Errorf("error preparing multi-search request: %w", err)
	}

	// Execute the multi-search request
	res, err := es.client.Msearch(
		bytes.NewReader(buffer.Bytes()),
		es.client.Msearch.WithIndex(index),
		es.client.Msearch.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing multi-search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error in multi-search response: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing multi-search response: %w", err)
	}

	responses := result["responses"].([]interface{})
	results := make([][]json.RawMessage, len(responses))

	for i, response := range responses {
		resp := response.(map[string]interface{})
		if resp["error"] != nil {
			return nil, fmt.Errorf("error in search response %d: %v", i, resp["error"])
		}

		hits := resp["hits"].(map[string]interface{})["hits"].([]interface{})
		documents := make([]json.RawMessage, len(hits))
		for j, hit := range hits {
			source := hit.(map[string]interface{})["_source"]
			docBytes, err := json.Marshal(source)
			if err != nil {
				return nil, fmt.Errorf("error marshaling document: %w", err)
			}
			documents[j] = docBytes
		}
		results[i] = documents
	}

	return results, nil
}

// DeleteIndex deletes an index from Elasticsearch
func (es *ElasticsearchClient) DeleteIndex(ctx context.Context, index string) error {
	res, err := es.client.Indices.Delete([]string{index})
	if err != nil {
		return fmt.Errorf("error deleting index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	return nil
}
