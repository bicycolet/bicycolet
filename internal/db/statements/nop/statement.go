package nop

import (
	"github.com/bicycolet/bicycolet/internal/db/statements/statement"
)

// Registrar defines a nop implementation of statement registry
type Registrar struct{}

// New creates a new registrar.
func New() *Registrar {
	return &Registrar{}
}

// Create returns a statement if it's found, or will prepare on for use.
func (Registrar) Create(statement string) (statement.Statement, error) {
	return nil, nil
}
