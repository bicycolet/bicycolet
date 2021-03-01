package queries

import (
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	"github.com/bicycolet/bicycolet/internal/db/queries/sqlite"
	"github.com/pkg/errors"
)

// Query performs queries against the underlying database.
type Query struct {
	statements query.Statements
}

// StatementType defines the statement we want to use for resilience.
type StatementType int

const (
	// SQLite statement that describes what statement to use.
	SQLite StatementType = iota

	// TODO (stickupkid): Implement postgres queries.
)

// New creates a new encoding statement based on the type.
func New(t StatementType, options ...Option) (query.Queries, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	switch t {
	case SQLite:
		return &Query{
			statements: sqlite.Statements{},
		}, nil
	default:
		return nil, errors.Errorf("invalid type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct{}

// Create a options instance with default values.
func newOptions() *options {
	return &options{}
}
