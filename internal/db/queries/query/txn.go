package query

// Txn is an in-progress database transaction.
// A transaction must end with a call to Commit or Rollback.
type Txn interface {
	// Query executes a query that returns rows, typically a SELECT.
	Query(query string, args ...interface{}) (Rows, error)
}
