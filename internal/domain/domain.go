package domain

import (
	"encoding/json"
	"time"
)

// Query represents the query filters for retrieving the logs
type Query struct {
	Entries        int    // number of lines to fetch
	Query          string // optional term query
	Reverse        bool
	AfterDateTime  *time.Time
	BeforeDateTime *time.Time
	IndexPattern   string
	Refresh        time.Duration
	Format         string
	FormatFields   []string
	TimestampField string
	ShowTime       bool
}

// LogEntry represents a log entry fetched from the database
type LogEntry struct {
	ID      string
	Message *json.RawMessage
}
