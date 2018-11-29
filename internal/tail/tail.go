package tail

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/pmdcosta/elklogs/internal/domain"
	"github.com/sirupsen/logrus"
)

// Connector abstracts the database connection
type Connector interface {
	Close() error
	GetIndexNames(ctx context.Context) ([]string, error)
	ExecuteQuery(ctx context.Context, indices []string, timestampField string, order bool, query string, entries int) ([]*domain.LogEntry, error)
}

// Tail is a structure that holds data necessary to perform tailing
type Tail struct {
	logger    *logrus.Entry
	connector Connector

	lastID string
}

// New creates a new Tail
func New(logger *logrus.Entry, connector Connector) *Tail {
	t := &Tail{
		logger:    logger,
		connector: connector,
	}

	return t
}

// Start starts tailing logs
func (t *Tail) Start(query *domain.Query) error {
	// get cluster indices
	indices, err := t.connector.GetIndexNames(context.Background())
	if err != nil {
		return errors.Wrap(err, "could not fetch available indices")
	}
	t.logger.WithFields(logrus.Fields{"indices": indices}).Debug("indices fetched")

	// filter indices based on query date filters
	indices, err = FilterIndex(indices, query.IndexPattern, query.AfterDateTime, query.BeforeDateTime)
	if err != nil {
		return errors.Wrap(err, "could not filter indices")
	}
	t.logger.WithFields(logrus.Fields{"indices": indices}).Debug("indices filtered")

	// execute
	if err = t.loop(query, indices); err != nil {
		return err
	}

	// tail the logs
	for query.Refresh != 0 {
		// refresh timer
		time.Sleep(query.Refresh)
		if err = t.loop(query, indices); err != nil {
			return err
		}
	}

	return nil
}

// loop retrieves logs from the database, processes them and prints them
func (t *Tail) loop(query *domain.Query, indices []string) error {
	// retrieve logs from the host
	logs, err := t.connector.ExecuteQuery(context.Background(), indices, timestampField, false, query.Query, query.Entries)
	if err != nil {
		return errors.Wrap(err, "could not fetch logs")
	}
	t.logger.WithFields(logrus.Fields{"indices": indices, "query": query.Query, "entries": query.Entries, "logs": len(logs)}).Debug("logs fetched")

	// process logs
	entries, err := t.processLogs(query, logs)
	if err != nil {
		return errors.Wrap(err, "could not process logs")
	}
	t.logger.WithFields(logrus.Fields{"logs": len(entries)}).Debug("logs processed")

	// print logs
	printLogs(entries, query.Reverse)
	return nil
}

// processLogs processes json messages and returns the log entries according to the provided format
func (t *Tail) processLogs(query *domain.Query, logs []*domain.LogEntry) ([]string, error) {
	entries := make([]string, 0, len(logs))
	for _, log := range logs {
		if log.ID == t.lastID {
			return entries, nil
		}

		s, err := processEntry(log, query.ShowTime, query.Format, query.FormatFields)
		if err != nil {
			return nil, err
		}
		entries = append(entries, s)
	}

	// we need to keep track of the ID of the last message to remove duplicates between loops
	if t.lastID == "" {
		t.lastID = logs[0].ID
	}

	return entries, nil
}
