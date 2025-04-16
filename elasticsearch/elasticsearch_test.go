package elasticsearch

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupTestIndices(t *testing.T, client *ElasticsearchClient) {
	indices := []string{"test-index", "test-search-index", "empty-index"}
	for _, index := range indices {
		res, err := client.client.Indices.Delete([]string{index})
		if err != nil {
			t.Logf("Warning: Failed to delete index %s: %v", index, err)
			continue
		}
		defer res.Body.Close()
	}
}

func TestElasticsearchClient_IndexDocument(t *testing.T) {
	client, err := NewElasticsearchClient("localhost", "10200")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Cleanup before test
	cleanupTestIndices(t, client)

	// Ensure cleanup after test
	defer cleanupTestIndices(t, client)

	tests := []struct {
		name    string
		index   string
		doc     interface{}
		wantErr bool
	}{
		{
			name:  "successful index",
			index: "test-index",
			doc: map[string]interface{}{
				"title":   "Test Document",
				"content": "This is a test document",
			},
			wantErr: false,
		},
		{
			name:    "empty index name",
			index:   "",
			doc:     map[string]interface{}{"test": "data"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.IndexDocument(context.Background(), tt.index, tt.doc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestElasticsearchClient_SearchDocuments(t *testing.T) {
	client, err := NewElasticsearchClient("localhost", "10200")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Cleanup before test
	cleanupTestIndices(t, client)

	// Ensure cleanup after test
	defer cleanupTestIndices(t, client)

	// First index a test document
	testDoc := map[string]interface{}{
		"title":   "Test Search Document",
		"content": "This is a test document for search",
	}
	err = client.IndexDocument(context.Background(), "test-search-index", testDoc)
	require.NoError(t, err)

	// Wait for the document to be indexed
	time.Sleep(1 * time.Second)

	tests := []struct {
		name    string
		index   string
		query   map[string]interface{}
		wantErr bool
	}{
		{
			name:  "successful search",
			index: "test-search-index",
			query: map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "Test",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty index name",
			index:   "",
			query:   map[string]interface{}{"query": map[string]interface{}{"match_all": map[string]interface{}{}}},
			wantErr: false,
		},
		{
			name:  "empty index search",
			index: "empty-index",
			query: map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := client.SearchDocuments(context.Background(), tt.index, tt.query)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, docs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, docs)
				if len(docs) > 0 {
					var result map[string]interface{}
					err := json.Unmarshal(docs[0], &result)
					assert.NoError(t, err)
					assert.NotEmpty(t, result)
				}
			}
		})
	}
}

func TestNewElasticsearchClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful client creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewElasticsearchClient("localhost", "10200")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}
