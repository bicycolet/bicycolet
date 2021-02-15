package query

// Table defines a table location for a given query.
type Table string

// Where clause for a given query.
type Where string

// Statements defines generic queries for a given underlying query engine.
type Statements interface {
	// Count returns the number of rows in the given table.
	Count(Table, Where) Query
}
