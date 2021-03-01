package query

// Txn is an in-progress database transaction.
// A transaction must end with a call to Commit or Rollback.
type Txn interface {
	// Query executes a query that returns rows, typically a SELECT.
	Query(string, ...interface{}) (Rows, error)

	// Exec executes a query that returns a result, typically an INSERT.
	Exec(string, ...interface{}) (Result, error)
}
