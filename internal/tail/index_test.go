package tail_test

import (
	"testing"
	"time"

	"github.com/pmdcosta/elklogs/internal/tail"
	"github.com/stretchr/testify/assert"
)

const pattern = "logstash-[0-9].*"

func TestFilterIndex(t *testing.T) {
	var indices = []string{"logstash-2018.11.03", "logstash-2018.10.10", "logstash-2018.10.09", "logstash-2018.10.11"}

	r, err := tail.FilterIndex(indices, pattern, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, []string{"logstash-2018.11.03"}, r)
}

func TestFilterIndex_end(t *testing.T) {
	var indices = []string{"logstash-2018.11.03", "logstash-2018.10.10", "logstash-2018.10.09", "logstash-2018.10.11"}

	end, err := time.Parse("2006-01-02T15:04", "2018-11-01T00:00")
	assert.Nil(t, err)
	r, err := tail.FilterIndex(indices, pattern, nil, &end)
	assert.Nil(t, err)
	assert.Equal(t, []string{"logstash-2018.10.10", "logstash-2018.10.09", "logstash-2018.10.11"}, r)
}

func TestFilterIndex_start(t *testing.T) {
	var indices = []string{"logstash-2018.11.03", "logstash-2018.10.10", "logstash-2018.10.09", "logstash-2018.10.11"}

	start, err := time.Parse("2006-01-02T15:04", "2018-10-10T00:00")
	assert.Nil(t, err)
	r, err := tail.FilterIndex(indices, pattern, &start, nil)
	assert.Nil(t, err)
	assert.Equal(t, []string{"logstash-2018.11.03", "logstash-2018.10.10", "logstash-2018.10.11"}, r)
}
