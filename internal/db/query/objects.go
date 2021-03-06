package query

import (
	"fmt"
	"strings"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/pkg/errors"
)

// Dest is a function that is expected to return the objects to pass to the
// 'dest' argument of sql.Rows.Scan(). It is invoked by SelectObjects once per
// yielded row, and it will be passed the index of the row being scanned.
type Dest func(i int) []interface{}

// SelectObjects executes a statement which must yield rows with a specific
// columns schema. It invokes the given Dest hook for each yielded row.
func SelectObjects(tx database.Tx, dest Dest, query string, args ...interface{}) error {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(dest(i)...); err != nil {
			return errors.WithStack(err)
		}
	}
	err = rows.Err()
	return errors.WithStack(err)
}

// UpsertObject inserts or replaces a new row with the given column values, to
// the given table using columns order. For example:
//
// UpsertObject(tx, "cars", []string{"id", "brand"}, []interface{}{1, "ferrari"})
//
// The number of elements in 'columns' must match the one in 'values'.
func UpsertObject(tx database.Tx, table string, columns []string, values []interface{}) (int64, error) {
	n := len(columns)
	if n == 0 {
		return -1, errors.Errorf("columns length is zero")
	}
	if n != len(values) {
		return -1, errors.Errorf("columns length does not match values length")
	}

	stmt := fmt.Sprintf(
		"INSERT OR REPLACE INTO %s (%s) VALUES %s",
		table, strings.Join(columns, ", "), Params(n))
	result, err := tx.Exec(stmt, values...)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return id, nil
}

// DeleteObject removes the row identified by the given ID. The given table
// must have a primary key column called 'id'.
//
// It returns a flag indicating if a matching row was actually found and
// deleted or not.
func DeleteObject(tx database.Tx, table string, id int64) (bool, error) {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id=?", table)
	result, err := tx.Exec(stmt, id)
	if err != nil {
		return false, errors.WithStack(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, errors.WithStack(err)
	}
	if n > 1 {
		return true, errors.Errorf("more than one row was deleted")
	}
	return n == 1, nil
}
