package sql

import (
	"sync"

	"github.com/bicycolet/bicycolet/internal/db/statements/statement"
)

// Registrar holds all the statements within a store, that are used to prepare
// them for use.
//
// The store can either register them ahead of time or just in time.
type Registrar struct {
	mutex      sync.RWMutex
	hasher     statement.Hasher
	preparer   statement.Preparer
	statements map[string]Statement
}

// New creates a new registry to hold all the statements.
func New(preparer statement.Preparer, hasher statement.Hasher) *Registrar {
	return &Registrar{
		hasher:     hasher,
		preparer:   preparer,
		statements: make(map[string]Statement),
	}
}

// Create returns a statement if it's found, or will prepare on for use.
func (r *Registrar) Create(statement string) (statement.Statement, error) {
	hash := r.hasher.Hash(statement)

	r.mutex.RLock()
	stmt, ok := r.statements[hash]
	r.mutex.RUnlock()

	if ok {
		return stmt.Statement, nil
	}

	prepared, err := r.preparer.Prepare(statement)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.statements[hash] = Statement{
		Raw:       statement,
		Statement: prepared,
	}

	return prepared, nil
}

// Statement holds the raw statement along with the prepared statement.
type Statement struct {
	Raw       string
	Statement statement.Statement
}
