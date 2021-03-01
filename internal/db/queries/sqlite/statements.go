package sqlite

import (
	"fmt"
	"strings"

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

// Delete removes the row identified by the given ID. The given table
// must have a primary key column called 'id'.
func (Statements) Delete(table query.Table) query.Query {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id=?", table)

	return queryRunner{
		stmt: stmt,
	}
}

// UpsertObject inserts or replaces a new row with the given column values, to
// the given table using columns order.
func (s Statements) UpsertObject(table query.Table, columns []string) query.Query {
	stmt := fmt.Sprintf("INSERT OR REPLACE INTO %s (%s) VALUES %s",
		table, strings.Join(columns, ", "), s.Params(len(columns)))

	return queryRunner{
		stmt: stmt,
	}
}

// Params returns a parameters expression with the given number of '?'
// placeholders. E.g. Params(2) -> "(?, ?)". Useful for IN and VALUES
// expressions.
func (Statements) Params(n int) string {
	tokens := make([]string, n)
	for i := 0; i < n; i++ {
		tokens[i] = "?"
	}
	return fmt.Sprintf("(%s)", strings.Join(tokens, ", "))
}

type queryRunner struct {
	stmt string
}

func (q queryRunner) Run(txn query.Txn, args ...interface{}) (query.Rows, error) {
	return txn.Query(q.stmt, args...)
}

func (q queryRunner) Exec(txn query.Txn, args ...interface{}) (query.Result, error) {
	return txn.Exec(q.stmt, args...)
}
