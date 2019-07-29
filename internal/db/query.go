package db

import (
	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/query"
)

// ObjectsQuery defines queries to the database for generic object queries
type ObjectsQuery interface {
	// SelectObjects executes a statement which must yield rows with a specific
	// columns schema. It invokes the given Dest hook for each yielded row.
	SelectObjects(database.Tx, query.Dest, string, ...interface{}) error

	// UpsertObject inserts or replaces a new row with the given column values,
	// to the given table using columns order. For example:
	//
	// UpsertObject(tx, "cars", []string{"id", "brand"}, []interface{}{1, "ferrari"})
	//
	// The number of elements in 'columns' must match the one in 'values'.
	UpsertObject(database.Tx, string, []string, []interface{}) (int64, error)

	// DeleteObject removes the row identified by the given ID. The given table
	// must have a primary key column called 'id'.
	//
	// It returns a flag indicating if a matching row was actually found and
	// deleted or not.
	DeleteObject(database.Tx, string, int64) (bool, error)
}

// StringsQuery defines queries to the database for string queries
type StringsQuery interface {

	// SelectStrings executes a statement which must yield rows with a single
	// string column. It returns the list of column values.
	SelectStrings(database.Tx, string, ...interface{}) ([]string, error)
}

// CountQuery defines queries to the database for count queries
type CountQuery interface {

	// Count returns the number of rows in the given table.
	Count(database.Tx, string, string, ...interface{}) (int, error)
}

// Query defines different queries for accessing the database
type Query interface {
	ObjectsQuery
	StringsQuery
	CountQuery
}

// Transaction defines a method for executing transactions over the
// database
type Transaction interface {
	// Transaction executes the given function within a database transaction.
	Transaction(database.DB, func(database.Tx) error) error
}
