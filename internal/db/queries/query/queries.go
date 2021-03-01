package query

// Queries holds all the functions for a query.
type Queries interface {
	// Count returns the number of rows in the given table.
	Count(Txn, Table, Where, ...interface{}) (int, error)
}
