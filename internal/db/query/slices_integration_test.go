// +build integration

package query_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/query"
	"github.com/bicycolet/bicycolet/internal/testing"
)

var testStringsErrorCases = []struct {
	query string
	err   string
}{
	{"garbage", "near \"garbage\": syntax error"},
	{"SELECT id, name FROM test", "query yields 2 columns, not 1"},
	{"SELECT id FROM test", "query yields \"INTEGER\" column, not \"TEXT\""},
}

var testIntegersErrorCases = []struct {
	query string
	err   string
}{
	{"garbage", "near \"garbage\": syntax error"},
	{"SELECT id, name FROM test", "query yields 2 columns, not 1"},
	{"SELECT name FROM test", "query yields \"TEXT\" column, not \"INTEGER\""},
}

// Exercise possible failure modes.
func TestStrings_Error(t *testing.T) {
	for _, c := range testStringsErrorCases {
		t.Run(c.query, func(t *testing.T) {
			tx, close := newTxForSlices(t)
			defer close()

			values, err := query.SelectStrings(tx, c.query)
			if expected, actual := c.err, err.Error(); expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
			if values != nil {
				t.Errorf("expected values to be nil: %v", values)
			}
		})
	}
}

// All values yield by the query are returned.
func TestStrings(t *testing.T) {
	tx, close := newTxForSlices(t)
	defer close()

	values, err := query.SelectStrings(tx, "SELECT name FROM test ORDER BY name")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := []string{"bar", "foo"}, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Exercise possible failure modes.
func TestIntegers_Error(t *testing.T) {
	for _, c := range testIntegersErrorCases {
		t.Run(c.query, func(t *testing.T) {
			tx, close := newTxForSlices(t)
			defer close()

			values, err := query.SelectIntegers(tx, c.query)
			if expected, actual := c.err, err.Error(); expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
			if values != nil {
				t.Errorf("expected values to be nil: %v", values)
			}
		})
	}
}

// All values yield by the query are returned.
func TestIntegers(t *testing.T) {
	tx, close := newTxForSlices(t)
	defer close()

	values, err := query.SelectIntegers(tx, "SELECT id FROM test ORDER BY id")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := []int{0, 1}, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Insert new rows in bulk.
func TestInsertStrings_Bulk(t *testing.T) {
	tx, close := newTxForSlices(t)
	defer close()

	err := query.InsertStrings(tx, "INSERT INTO test(name) VALUES %s", []string{"xx", "yy"})
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	values, err := query.SelectStrings(tx, "SELECT name FROM test ORDER BY name DESC LIMIT 2")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := []string{"yy", "xx"}, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Return a new transaction against an in-memory postgres database with a single
// test table populated with a few rows.
func newTxForSlices(t *testing.T) (database.Tx, func()) {
	connInfo, err := testing.ConnectionInfo()
	if err != nil {
		t.Fatalf("expected err to be nil: %v", err)
	}
	db, err := sql.Open(database.DriverName(), connInfo.String())
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("INSERT INTO test VALUES (0, 'foo'), (1, 'bar')")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	tx, err := database.ShimTx(db.Begin())
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	return tx, func() {
		db.Close()
	}
}
