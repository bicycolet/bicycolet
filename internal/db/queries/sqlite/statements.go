package sqlite

import (
	"fmt"

	"github.com/bicycolet/bicycolet/internal/db/queries/query"
)

// Statements holds the statements for a given sqlite database.
type Statements struct{}

// Count defines a generic count statement.
func (Statements) Count(table query.Table, where query.Where) query.Query {
	stmt := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if where != "" {
		stmt += fmt.Sprintf(" WHERE %s", where)
	}

	return queryRunner{
		stmt: stmt,
	}
}

type queryRunner struct {
	stmt string
}

func (q queryRunner) Run(txn query.Txn, args ...interface{}) (query.Rows, error) {
	return txn.Query(q.stmt, args...)
}
