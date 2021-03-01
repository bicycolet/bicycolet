package queries_test

import (
	"database/sql"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/queries"
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	_ "github.com/mattn/go-sqlite3"
)

// Exercise possible failure modes.
func TestSelectObjects_Error(t *testing.T) {
	cases := []struct {
		dest  query.Dest
		query string
		err   string
	}{
		{
			func(int) []interface{} { return nil },
			"garbage",
			"near \"garbage\": syntax error",
		},
		{
			func(int) []interface{} { return make([]interface{}, 1) },
			"SELECT id, name FROM test",
			"sql: expected 2 destination arguments in Scan, not 1",
		},
	}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			tx := newTxForObjects(t)

			queries, err := queries.New(queries.SQLite)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}
			err = queries.SelectObjects(tx, c.dest, c.query)
			if expected, actual := c.err, err.Error(); expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

// Scan rows yielded by the query.
func TestSelectObjects_Success(t *testing.T) {
	tx := newTxForObjects(t)
	objects := make([]struct {
		ID   int
		Name string
	}, 1)
	object := objects[0]

	dest := func(i int) []interface{} {
		if expected, actual := 0, i; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		return []interface{}{
			&object.ID,
			&object.Name,
		}
	}

	stmt := "SELECT id, name FROM test WHERE name=?"

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	err = queries.SelectObjects(tx, dest, stmt, "bar")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	if expected, actual := 1, object.ID; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "bar", object.Name; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Exercise possible failure modes.
func TestUpsertObject_Error(t *testing.T) {
	cases := []struct {
		columns []string
		values  []interface{}
		err     string
	}{
		{
			[]string{},
			[]interface{}{},
			"columns length is zero",
		},
		{
			[]string{"id"},
			[]interface{}{2, "egg"},
			"columns length does not match values length",
		},
	}
	for _, c := range cases {
		t.Run(c.err, func(t *testing.T) {
			tx := newTxForObjects(t)

			queries, err := queries.New(queries.SQLite)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}
			err = queries.UpsertObject(tx, "foo", c.columns, c.values)
			if expected, actual := c.err, err.Error(); expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

// Insert a new row.
func TestUpsertObject_Insert(t *testing.T) {
	tx := newTxForObjects(t)

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	err = queries.UpsertObject(tx, "test", []string{"name"}, []interface{}{"egg"})
	if expected, actual := true, err == nil; expected != actual {
		t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
	}

	objects := make([]struct {
		ID   int
		Name string
	}, 1)
	object := objects[0]

	dest := func(i int) []interface{} {
		if expected, actual := 0, i; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		return []interface{}{
			&object.ID,
			&object.Name,
		}
	}

	stmt := "SELECT id, name FROM test WHERE name=?"

	err = queries.SelectObjects(tx, dest, stmt, "egg")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	if expected, actual := 2, object.ID; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "egg", object.Name; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Update an existing row.
func TestUpsertObject_Update(t *testing.T) {
	tx := newTxForObjects(t)

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	err = queries.UpsertObject(tx, "test", []string{"id", "name"}, []interface{}{1, "egg"})
	if expected, actual := true, err == nil; expected != actual {
		t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
	}

	objects := make([]struct {
		ID   int
		Name string
	}, 1)
	object := objects[0]

	dest := func(i int) []interface{} {
		if expected, actual := 0, i; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		return []interface{}{
			&object.ID,
			&object.Name,
		}
	}

	stmt := "SELECT id, name FROM test WHERE name=?"
	err = queries.SelectObjects(tx, dest, stmt, "egg")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	if expected, actual := 1, object.ID; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "egg", object.Name; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Exercise possible failure modes.
func TestDeleteObject_Error(t *testing.T) {
	tx := newTxForObjects(t)

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	deleted, err := queries.DeleteObject(tx, "foo", 1)
	if expected, actual := "no such table: foo", err.Error(); expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := false, deleted; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// If an row was actually deleted, the returned flag is true.
func TestDeleteObject_Deleted(t *testing.T) {
	tx := newTxForObjects(t)

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	deleted, err := queries.DeleteObject(tx, "test", 1)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := true, deleted; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// If no row was actually deleted, the returned flag is false.
func TestDeleteObject_NotDeleted(t *testing.T) {
	tx := newTxForObjects(t)

	queries, err := queries.New(queries.SQLite)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	deleted, err := queries.DeleteObject(tx, "test", 1000)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := false, deleted; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Return a new transaction against an in-memory SQLite database with a single
// test table populated with a few rows for testing object-related queries.
func newTxForObjects(t *testing.T) query.Txn {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	_, err = db.Exec("INSERT INTO test VALUES (0, 'foo'), (1, 'bar')")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	txn, err := db.Begin()
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	return shimTxn{txn: txn}
}
