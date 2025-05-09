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
	res, err := client.client.Indices.Delete([]string{"_all"})
	if err != nil {
		t.Logf("Warning: Failed to delete all indices: %v", err)
		return
	}
	defer res.Body.Close()
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
		query   *ESQuery
		wantErr bool
	}{
		{
			name: "successful search",
			query: CreateESQuery("test-search-index", map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "Test",
					},
				},
			}),
			wantErr: false,
		},
		{
			name: "empty index name",
			query: CreateESQuery("", map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			}),
			wantErr: true,
		},
		{
			name: "empty index search",
			query: CreateESQuery("empty-index", map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := client.SearchDocuments(context.Background(), tt.query)
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

func TestElasticsearchClient_DeleteIndex(t *testing.T) {
	client, err := NewElasticsearchClient("localhost", "10200")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Cleanup before test
	cleanupTestIndices(t, client)

	// Ensure cleanup after test
	defer cleanupTestIndices(t, client)

	// First create a test index by indexing a document
	testDoc := map[string]interface{}{
		"title":   "Test Document",
		"content": "This is a test document",
	}
	err = client.IndexDocument(context.Background(), "test-delete-index", testDoc)
	require.NoError(t, err)

	// Wait for the document to be indexed
	time.Sleep(1 * time.Second)

	tests := []struct {
		name    string
		index   string
		wantErr bool
	}{
		{
			name:    "successful index deletion",
			index:   "test-delete-index",
			wantErr: false,
		},
		{
			name:    "delete non-existent index",
			index:   "non-existent-index",
			wantErr: true,
		},
		{
			name:    "empty index name",
			index:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteIndex(context.Background(), tt.index)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestElasticsearchClient_Msearch(t *testing.T) {
	client, err := NewElasticsearchClient("localhost", "10200")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Cleanup before test
	cleanupTestIndices(t, client)

	// First index test documents in different indices
	testDocs := []struct {
		index string
		doc   map[string]interface{}
	}{
		{
			index: "test-msearch-index-1",
			doc: map[string]interface{}{
				"title":   "Test Document 1",
				"content": "This is the first test document",
			},
		},
		{
			index: "test-msearch-index-2",
			doc: map[string]interface{}{
				"title":   "Test Document 2",
				"content": "This is the second test document",
			},
		},
	}

	for _, testDoc := range testDocs {
		err := client.IndexDocument(context.Background(), testDoc.index, testDoc.doc)
		require.NoError(t, err)
	}

	// Wait for the documents to be indexed
	time.Sleep(1 * time.Second)

	// Create a multi-search query that searches across different indices
	mquery := CreateMQuery()
	mquery.AddQuery(CreateESQuery("test-msearch-index-1", map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Document 1",
			},
		},
	}))
	mquery.AddQuery(CreateESQuery("test-msearch-index-2", map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Document 2",
			},
		},
	}))

	// Execute the multi-search
	results, err := client.SearchDocumentsWithQuery(context.Background(), "test-msearch-index-1", mquery)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Verify first query results from first index
	assert.Len(t, results[0], 1)
	var doc1 map[string]interface{}
	err = json.Unmarshal(results[0][0], &doc1)
	assert.NoError(t, err)
	assert.Equal(t, "Test Document 1", doc1["title"])

	// Verify second query results from second index
	assert.Len(t, results[1], 1)
	var doc2 map[string]interface{}
	err = json.Unmarshal(results[1][0], &doc2)
	assert.NoError(t, err)
	assert.Equal(t, "Test Document 2", doc2["title"])
}
