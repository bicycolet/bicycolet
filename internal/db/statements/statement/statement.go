package statement

import "context"

// Statement is a prepared statement.
type Statement interface {
	// ExecContext executes a prepared statement with the given arguments and
	// returns a Result summarizing the effect of the statement.
	ExecContext(context.Context, ...interface{}) (Result, error)
}

// A Result summarizes an executed statement command.
type Result interface {
	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected() (int64, error)
}
