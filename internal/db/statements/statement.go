package statements

import (
	"crypto/sha256"
	"fmt"

	"github.com/bicycolet/bicycolet/internal/db/statements/nop"
	"github.com/bicycolet/bicycolet/internal/db/statements/sql"
	"github.com/bicycolet/bicycolet/internal/db/statements/statement"
	"github.com/pkg/errors"
)

// StatementType defines the statement we want to use for resilience.
type StatementType int

const (
	// SQL statement that describes what statement to use.
	SQL StatementType = iota

	// Nop statement that describes what statement to use.
	Nop
)

// New creates a new encoding statement based on the type.
func New(t StatementType, prep statement.Preparer, options ...Option) (statement.Registrar, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	switch t {
	case SQL:
		return sql.New(prep, opts.hasher), nil
	case Nop:
		return nop.New(), nil
	default:
		return nil, errors.Errorf("invalid type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	hasher statement.Hasher
}

// WithHasher sets the hasher for statements.
func WithHasher(hasher statement.Hasher) Option {
	return func(options *options) {
		options.hasher = hasher
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		hasher: stdlibHasher{},
	}
}

type stdlibHasher struct{}

func (stdlibHasher) Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:])
}
