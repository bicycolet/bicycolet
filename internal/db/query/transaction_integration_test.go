// +build integration

package query_test

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/query"
	"github.com/pkg/errors"
)

// Any error happening when beginning the transaction will be propagated.
func TestTransaction_BeginError(t *testing.T) {
	db := newDB(t)
	db.Close()

	err := query.Transaction(db, func(database.Tx) error { return nil })
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
	if expected, actual := "failed to begin transaction", err.Error(); !strings.Contains(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Any error happening when in the transaction function will cause a rollback.
func TestTransaction_FunctionError(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	err := query.Transaction(db, func(tx database.Tx) error {
		_, err := tx.Exec("CREATE TABLE test (id INTEGER)")
		if err != nil {
			t.Errorf("expected err to be nil: %v", err)
		}
		return errors.Errorf("boom")

	})
	if expected, actual := "boom", err.Error(); expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	tables, err := query.SelectStrings(tx, "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE'")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := "test", tables; contains(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Return a new in-memory SQLite database.
func newDB(t *testing.T) database.DB {
	db, err := sql.Open(database.DriverName(), connectionInfo())
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	return database.NewShimDB(db)
}

func contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

func connectionInfo() string {
	info := database.ConnectionInfo{
		Host:     "localhost",
		Port:     5435,
		User:     "postgres",
		Password: "postgres",
		DBName:   "test",
	}
	return info.String()
}
