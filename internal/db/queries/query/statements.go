package query

// Table defines a table location for a given query.
type Table string

// Where clause for a given query.
type Where string

// Statements defines generic queries for a given underlying query engine.
type Statements interface {
	// Count returns the number of rows in the given table.
	Count(Table, Where) Query

	// Delete removes the row identified by the given ID.
	Delete(Table) Query

	// UpsertObject inserts or replaces a new row with the given column values,
	// to the given table using columns order.
	UpsertObject(Table, []string) Query
}
