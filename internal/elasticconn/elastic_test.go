package elasticconn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pmdcosta/elklogs/internal/elasticconn"
	"github.com/stretchr/testify/assert"
)

type Elastic struct {
	*elasticconn.Elastic
}

// MustCreateConnector returns a new connector for testing
func MustCreateConnector(t *testing.T) *Elastic {
	e, err := elasticconn.New("https://127.0.0.1", "", "")
	assert.Nil(t, err)
	c := &Elastic{
		e,
	}
	return c
}

func TestConnector_Connect(t *testing.T) {
	MustCreateConnector(t)
}

func TestConnector_GetIndexNames(t *testing.T) {
	c := MustCreateConnector(t)

	indices, err := c.GetIndexNames(context.Background())
	assert.Nil(t, err)
	fmt.Println(indices)
}

func TestConnector_ExecuteQuery(t *testing.T) {
	t.Run("empty query", testConnector_ExecuteQuery_empty)
	t.Run("query", testConnector_ExecuteQuery_query)
}

func testConnector_ExecuteQuery_empty(t *testing.T) {
	c := MustCreateConnector(t)

	r, err := c.ExecuteQuery(context.Background(), []string{"logstash-2018.11.28"}, "@timestamp", false, "", 10)
	assert.Nil(t, err)
	assert.Len(t, r, 10)
}

func testConnector_ExecuteQuery_query(t *testing.T) {
	c := MustCreateConnector(t)

	r, err := c.ExecuteQuery(context.Background(), []string{"logstash-2018.11.29"}, "@timestamp", false, "kubernetes.labels.app:test", 20)
	assert.Nil(t, err)
	assert.Len(t, r, 20)
}
