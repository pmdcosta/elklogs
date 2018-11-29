package tail_test

import (
	"encoding/json"
	"testing"

	"github.com/pmdcosta/elklogs/internal/tail"
	"github.com/stretchr/testify/assert"
)

func TestProcessLogs(t *testing.T) {
	a := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:34","test":"message\n"}`)
	b := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:35","test":"message2\n"}`)

	logs := []*json.RawMessage{&a, &b}
	format := "%@timestamp: %test"
	fields := []string{"%@timestamp", "%test"}

	r, err := tail.ProcessLogs(logs, false, format, fields)
	assert.Nil(t, err)
	assert.Equal(t, []string{"2018-11-29T04:51:34: message", "2018-11-29T04:51:35: message2"}, r)
}

func TestProcessLogs_empty(t *testing.T) {
	a := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:34","test":"message\n"}`)
	b := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:35","test":"message2\n"}`)

	logs := []*json.RawMessage{&a, &b}
	format := ""
	fields := []string{}

	r, err := tail.ProcessLogs(logs, false, format, fields)
	assert.Nil(t, err)
	assert.Equal(t, []string{`{"@timestamp":"2018-11-29T04:51:34","test":"message\n"}`, `{"@timestamp":"2018-11-29T04:51:35","test":"message2\n"}`}, r)
}

func TestProcessLogs_notfound(t *testing.T) {
	a := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:34","test":"message\n"}`)
	b := json.RawMessage(`{"@timestamp":"2018-11-29T04:51:35","test":"message2\n"}`)

	logs := []*json.RawMessage{&a, &b}
	format := "%@timestamp: %test :%stuff"
	fields := []string{"%@timestamp", "%test", "%stuff"}

	r, err := tail.ProcessLogs(logs, false, format, fields)
	assert.Nil(t, err)
	assert.Equal(t, []string{"2018-11-29T04:51:34: message :", "2018-11-29T04:51:35: message2 :"}, r)
}
