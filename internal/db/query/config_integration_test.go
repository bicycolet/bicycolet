// +build integration

package query_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/query"
	internaltesting "github.com/bicycolet/bicycolet/internal/testing"
)

func TestSelectConfig_Selects(t *testing.T) {
	table, tx, close := newTxForConfig(t)
	defer close()

	values, err := query.SelectConfig(tx, table, "")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	want := map[string]string{
		"foo": "x",
		"bar": "zz",
	}
	if expected, actual := want, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestSelectConfig_WithFilters(t *testing.T) {
	table, tx, close := newTxForConfig(t)
	defer close()

	values, err := query.SelectConfig(tx, table, "key=$1", "bar")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	want := map[string]string{
		"bar": "zz",
	}
	if expected, actual := want, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// New keys are added to the table.
func TestUpdateConfig_NewKeys(t *testing.T) {
	table, tx, close := newTxForConfig(t)
	defer close()

	values := map[string]string{"foo": "y"}
	err := query.UpdateConfig(tx, table, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	values, err = query.SelectConfig(tx, table, "")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	want := map[string]string{
		"foo": "y",
		"bar": "zz",
	}
	if expected, actual := want, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Unset keys are deleted from the table.
func TestDeleteConfig_Delete(t *testing.T) {
	table, tx, close := newTxForConfig(t)
	defer close()

	values := map[string]string{"foo": ""}

	err := query.UpdateConfig(tx, table, values)

	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	values, err = query.SelectConfig(tx, table, "")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	want := map[string]string{
		"bar": "zz",
	}
	if expected, actual := want, values; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Return a new transaction against an in-memory postgres database with a single
// test table populated with a few rows.
func newTxForConfig(t *testing.T) (string, database.Tx, func()) {
	connInfo, err := internaltesting.ConnectionInfo()
	if err != nil {
		t.Fatalf("expected err to be nil: %v", err)
	}
	db, err := sql.Open(database.DriverName(), connInfo.String())
	if err != nil {
		t.Fatalf("expected err to be nil: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("expected err to be nil: %v", err)
	}

	table := internaltesting.RandomTableName()

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %q (key TEXT NOT NULL, value TEXT)", table))
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO %q VALUES ('foo', 'x'), ('bar', 'zz')", table))
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	tx, err := database.ShimTx(db.Begin())
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	return table, tx, func() {
		db.Close()
	}
}
