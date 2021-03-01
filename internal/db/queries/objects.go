package queries

import (
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	"github.com/pkg/errors"
)

// SelectObjects executes a statement which must yield rows with a specific
// columns schema. It invokes the given Dest hook for each yielded row.
func (q *Query) SelectObjects(tx query.Txn, dest query.Dest, query string, args ...interface{}) error {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = rows.Close() }()

	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(dest(i)...); err != nil {
			return errors.WithStack(err)
		}
	}
	return errors.WithStack(rows.Err())
}

// UpsertObject inserts or replaces a new row with the given column values, to
// the given table using columns order. For example:
//
// UpsertObject(tx, "cars", []string{"id", "brand"}, []interface{}{1, "ferrari"})
//
// The number of elements in 'columns' must match the one in 'values'.
func (q *Query) UpsertObject(tx query.Txn, table query.Table, columns []string, values []interface{}) error {
	n := len(columns)
	if n == 0 {
		return errors.Errorf("columns length is zero")
	} else if n != len(values) {
		return errors.Errorf("columns length does not match values length")
	}

	result, err := q.statements.UpsertObject(table, columns).Exec(tx, values...)
	if err != nil {
		return errors.WithStack(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	if affected != 1 {
		return errors.Errorf("nothing changed")
	}
	return nil
}

// DeleteObject removes the row identified by the given ID. The given table
// must have a primary key column called 'id'.
//
// It returns a flag indicating if a matching row was actually found and
// deleted or not.
func (q *Query) DeleteObject(tx query.Txn, table query.Table, id int64) (bool, error) {
	result, err := q.statements.Delete(table).Exec(tx, id)
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
