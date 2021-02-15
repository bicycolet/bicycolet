package statement

// Preparer defines a way to prepare a statement for using.
type Preparer interface {
	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	Prepare(string) (Statement, error)
}
