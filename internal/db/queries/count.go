package queries

import (
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	"github.com/pkg/errors"
)

// Count returns the number of rows in the given table.
func (q *Query) Count(tx query.Txn, table query.Table, where query.Where, args ...interface{}) (int, error) {
	rows, err := q.statements.Count(table, where).Run(tx, args...)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	defer func() { _ = rows.Close() }()

	// For sanity, make sure we read one and only one row.
	if !rows.Next() {
		return -1, errors.Errorf("no rows returned")
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return -1, errors.Errorf("failed to scan count column")
	}
	if rows.Next() {
		return -1, errors.Errorf("more than one row returned")
	}
	if err = rows.Err(); err != nil {
		return -1, errors.WithStack(err)
	}

	return count, nil
}
