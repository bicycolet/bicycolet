package queries_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/queries"
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	_ "github.com/mattn/go-sqlite3"
)

// Count returns the current number of rows.
func TestCount_Cases(t *testing.T) {
	cases := []struct {
		where query.Where
		args  []interface{}
		count int
	}{
		{
			"id=?",
			[]interface{}{999},
			0,
		},
		{
			"id=?",
			[]interface{}{1},
			1,
		},
		{
			"",
			[]interface{}{},
			2,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("count %d", c.count), func(t *testing.T) {
			tx := newTxForCount(t)
			queries, err := queries.New(queries.SQLite)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}
			count, err := queries.Count(tx, "test", c.where, c.args...)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}
			if expected, actual := c.count, count; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

// Return a new transaction against an in-memory SQLite database with a single
// test table and a few rows.
func newTxForCount(t *testing.T) query.Txn {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("CREATE TABLE test (id INTEGER)")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("INSERT INTO test VALUES (1), (2)")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	txn, err := db.Begin()
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	return shimTxn{txn: txn}
}

type shimTxn struct {
	txn *sql.Tx
}

// Query executes a query that returns rows, typically a SELECT.
func (s shimTxn) Query(query string, args ...interface{}) (query.Rows, error) {
	return s.txn.Query(query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
func (s shimTxn) Exec(query string, args ...interface{}) (query.Result, error) {
	return s.txn.Exec(query, args...)
}
