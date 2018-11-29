package elasticconn

import (
	"context"
	"time"

	"github.com/olivere/elastic"
	"github.com/pmdcosta/elklogs/internal/domain"
)

// Elastic manages the connection to the elasticsearch database
type Elastic struct {
	// elastic search client
	db     *elastic.Client
	config []elastic.ClientOptionFunc
}

// ElasticOptionFunc is a function that configures the Elastic Client.
type ElasticOptionFunc func(*Elastic) error

// OverrideElasticConfig can be used to override the default elastic search client configuration.
func OverrideElasticConfig(config []elastic.ClientOptionFunc) ElasticOptionFunc {
	return func(e *Elastic) error {
		e.config = config
		return nil
	}
}

// New creates a new elastic connector instance
func New(host string, user string, password string, options ...ElasticOptionFunc) (*Elastic, error) {
	e := &Elastic{}

	// default connection options
	defaultOptions := []elastic.ClientOptionFunc{
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeoutStartup(10 * time.Second),
		elastic.SetHealthcheckTimeout(2 * time.Second),
	}
	if user != "" {
		defaultOptions = append(defaultOptions, elastic.SetBasicAuth(user, password))
	}

	// run the optional options
	for _, option := range options {
		if err := option(e); err != nil {
			return nil, err
		}
	}

	// create the connection
	db, err := elastic.NewClient(defaultOptions...)
	if err != nil {
		return nil, err
	}

	// ping the database to make sure the connection was successful
	_, _, err = db.Ping(host).Do(context.Background())
	if err != nil {
		return nil, err
	}
	e.db = db

	return e, nil

}

// Close terminates the database connection
func (e *Elastic) Close() error {
	if e.db == nil {
		return nil
	}
	e.db = nil
	return nil
}

// GetIndexNames retrieves the names of all indices in the database
func (e Elastic) GetIndexNames(ctx context.Context) ([]string, error) {
	return e.db.IndexNames()
}

// ExecuteQuery executes the query against the database
func (e Elastic) ExecuteQuery(ctx context.Context, indices []string, timestampField string, order bool, query string, entries int) ([]*domain.LogEntry, error) {
	var q elastic.Query
	if query != "" {
		q = elastic.NewQueryStringQuery(query)
	} else {
		q = elastic.NewMatchAllQuery()
	}

	r, err := e.db.Search().Index(indices...).Sort(timestampField, order).Query(q).From(0).Size(entries).Do(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.LogEntry, 0, len(r.Hits.Hits))
	for _, i := range r.Hits.Hits {
		e := domain.LogEntry{
			ID:      i.Id,
			Message: i.Source,
		}
		result = append(result, &e)
	}

	return result, nil
}
