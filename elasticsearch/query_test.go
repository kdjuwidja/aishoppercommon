package elasticsearch

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateESQuery(t *testing.T) {
	tests := []struct {
		name      string
		index     string
		queryBody map[string]interface{}
		want      *ESQuery
		wantErr   bool
	}{
		{
			name:  "simple match query",
			index: "test-index",
			queryBody: map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "test",
					},
				},
			},
			want: &ESQuery{
				index: "test-index",
				query: map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "empty index",
			index: "",
			queryBody: map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "nil query body",
			index:     "test-index",
			queryBody: nil,
			want: &ESQuery{
				index: "test-index",
				query: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateESQuery(tt.index, tt.queryBody)
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestESQuery_appendBufferForMQuery(t *testing.T) {
	tests := []struct {
		name      string
		query     *ESQuery
		index     string
		wantBytes []byte
		wantErr   bool
	}{
		{
			name: "same index",
			query: &ESQuery{
				index: "test-index",
				query: map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				},
			},
			index: "test-index",
			wantBytes: []byte(`{ }
{"query":{"match":{"title":"test"}}}
`),
			wantErr: false,
		},
		{
			name: "different index",
			query: &ESQuery{
				index: "test-index-1",
				query: map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				},
			},
			index: "test-index-2",
			wantBytes: []byte(`{"index":"test-index-1"}
{"query":{"match":{"title":"test"}}}
`),
			wantErr: false,
		},
		{
			name: "empty query index",
			query: &ESQuery{
				index: "",
				query: map[string]interface{}{
					"query": map[string]interface{}{
						"match_all": map[string]interface{}{},
					},
				},
			},
			index:     "test-index",
			wantBytes: nil,
			wantErr:   true,
		},
		{
			name: "empty context index",
			query: &ESQuery{
				index: "test-index",
				query: map[string]interface{}{
					"query": map[string]interface{}{
						"match_all": map[string]interface{}{},
					},
				},
			},
			index:     "",
			wantBytes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.query.appendBufferForMQuery(tt.index, &buf)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, buf.Bytes())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, buf.Bytes())
			}
		})
	}
}

func TestCreateMQuery(t *testing.T) {
	mquery := CreateMQuery()
	assert.NotNil(t, mquery)
	assert.Empty(t, mquery.queries)
	assert.NotNil(t, mquery.queries) // ensure slice is initialized
}

func TestMQuery_AddQuery(t *testing.T) {
	tests := []struct {
		name      string
		mquery    *MultiESQuery
		query     *ESQuery
		wantLen   int
		wantQuery *ESQuery
	}{
		{
			name:   "add first query",
			mquery: CreateMQuery(),
			query: CreateESQuery("test-index-1", map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "test1",
					},
				},
			}),
			wantLen: 1,
			wantQuery: CreateESQuery("test-index-1", map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "test1",
					},
				},
			}),
		},
		{
			name: "add second query",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("test-index-1", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test1",
						},
					},
				}))
				return mq
			}(),
			query: CreateESQuery("test-index-2", map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "test2",
					},
				},
			}),
			wantLen: 2,
			wantQuery: CreateESQuery("test-index-2", map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"title": "test2",
					},
				},
			}),
		},
		{
			name:      "add nil query",
			mquery:    CreateMQuery(),
			query:     nil,
			wantLen:   0,
			wantQuery: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mquery.AddQuery(tt.query)
			assert.Len(t, tt.mquery.queries, tt.wantLen)
			if tt.wantQuery != nil {
				assert.Equal(t, tt.wantQuery, tt.mquery.queries[tt.wantLen-1])
			}
		})
	}
}

func TestMQuery_CreateMQueryBuffer(t *testing.T) {
	tests := []struct {
		name      string
		mquery    *MultiESQuery
		index     string
		wantBytes []byte
		wantErr   bool
	}{
		{
			name: "single query same index",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("test-index", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				}))
				return mq
			}(),
			index: "test-index",
			wantBytes: []byte(`{ }
{"query":{"match":{"title":"test"}}}
`),
			wantErr: false,
		},
		{
			name: "single query different index",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("test-index-1", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				}))
				return mq
			}(),
			index: "test-index-2",
			wantBytes: []byte(`{"index":"test-index-1"}
{"query":{"match":{"title":"test"}}}
`),
			wantErr: false,
		},
		{
			name: "multiple queries",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("test-index-1", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test1",
						},
					},
				}))
				mq.AddQuery(CreateESQuery("test-index-2", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test2",
						},
					},
				}))
				return mq
			}(),
			index: "test-index-3",
			wantBytes: []byte(`{"index":"test-index-1"}
{"query":{"match":{"title":"test1"}}}
{"index":"test-index-2"}
{"query":{"match":{"title":"test2"}}}
`),
			wantErr: false,
		},
		{
			name: "empty index",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("test-index", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				}))
				return mq
			}(),
			index:     "",
			wantBytes: nil,
			wantErr:   true,
		},
		{
			name: "empty query index",
			mquery: func() *MultiESQuery {
				mq := CreateMQuery()
				mq.AddQuery(CreateESQuery("", map[string]interface{}{
					"query": map[string]interface{}{
						"match": map[string]interface{}{
							"title": "test",
						},
					},
				}))
				return mq
			}(),
			index:     "test-index",
			wantBytes: nil,
			wantErr:   true,
		},
		{
			name: "empty mquery",
			mquery: func() *MultiESQuery {
				return CreateMQuery()
			}(),
			index:     "test-index",
			wantBytes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer, err := tt.mquery.createMQueryBuffer(tt.index)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, buffer)
			} else {
				assert.NoError(t, err)
				if tt.wantBytes == nil {
					assert.Nil(t, buffer)
				} else {
					assert.NotNil(t, buffer)
					assert.Equal(t, tt.wantBytes, buffer.Bytes())
				}
			}
		})
	}
}
