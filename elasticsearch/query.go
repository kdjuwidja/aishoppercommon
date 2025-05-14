package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ESQuery represents a single Elasticsearch query
type ESQuery struct {
	index string
	query map[string]interface{}
}

func CreateESQuery(index string, query map[string]interface{}) *ESQuery {
	if index == "" {
		return nil
	}
	return &ESQuery{
		index: index,
		query: query,
	}
}

func CreateESQueryStr(index string, queryStr string) *ESQuery {
	if index == "" {
		return nil
	}

	query := make(map[string]interface{})
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return nil
	}
	return &ESQuery{
		index: index,
		query: query,
	}
}

func (q *ESQuery) appendBufferForMQuery(index string, buffer *bytes.Buffer) error {
	if index == "" {
		return fmt.Errorf("index name cannot be empty")
	}

	if q.index == "" {
		return fmt.Errorf("index name cannot be empty")
	}

	queryBytes, err := json.Marshal(q.query)
	if err != nil {
		return err
	}

	// append index
	if q.index != index {
		// write index as context if index is different
		buffer.WriteString(fmt.Sprintf("{\"index\":\"%s\"}", q.index))
		buffer.WriteByte('\n')
	} else {
		// write empty context if index is the same
		buffer.WriteString("{ }")
		buffer.WriteByte('\n')
	}

	// append query
	buffer.Write(queryBytes)
	buffer.WriteByte('\n')

	return nil
}

// MultiESQuery represents a multi-query in Elasticsearch
type MultiESQuery struct {
	queries []*ESQuery
}

func CreateMQuery() *MultiESQuery {
	return &MultiESQuery{
		queries: []*ESQuery{},
	}
}

func (m *MultiESQuery) AddQuery(query *ESQuery) {
	if query != nil {
		m.queries = append(m.queries, query)
	}
}

func (m *MultiESQuery) createMQueryBuffer(index string) (*bytes.Buffer, error) {
	if len(m.queries) == 0 {
		return nil, fmt.Errorf("no queries to create multi-search buffer")
	}

	buffer := bytes.NewBuffer(nil)
	for _, query := range m.queries {
		if err := query.appendBufferForMQuery(index, buffer); err != nil {
			return nil, err
		}
	}

	return buffer, nil
}

func (m *MultiESQuery) PrintQuery(index string) {
	buffer, err := m.createMQueryBuffer(index)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buffer.String())
}
