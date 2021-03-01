package query

// Dest is a function that is expected to return the objects to pass to the
// 'dest' argument of sql.Rows.Scan(). It is invoked by SelectObjects once per
// yielded row, and it will be passed the index of the row being scanned.
type Dest func(i int) []interface{}

// Queries holds all the functions for a query.
type Queries interface {
	// Count returns the number of rows in the given table.
	Count(Txn, Table, Where, ...interface{}) (int, error)

	// SelectObjects executes a statement which must yield rows with a specific
	// columns schema. It invokes the given Dest hook for each yielded row.
	SelectObjects(Txn, Dest, string, ...interface{}) error

	// UpsertObject inserts or replaces a new row with the given column values, to
	// the given table using columns order. For example:
	//
	// UpsertObject(tx, "cars", []string{"id", "brand"}, []interface{}{1, "ferrari"})
	//
	// The number of elements in 'columns' must match the one in 'values'.
	UpsertObject(Txn, Table, []string, []interface{}) error

	// DeleteObject removes the row identified by the given ID. The given table
	// must have a primary key column called 'id'.
	//
	// It returns a flag indicating if a matching row was actually found and
	// deleted or not.
	DeleteObject(Txn, Table, int64) (bool, error)
}
