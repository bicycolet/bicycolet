package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

// Loading a config Map initializes it with the given values.
func TestNew(t *testing.T) {
	schema := Schema{
		"foo": {},
		"egg": NewBoolKey(),
	}

	cases := []struct {
		title  string
		values map[string]string // Initial values
		result map[string]string // Expected values after loading
	}{
		{
			title:  "plain load of regular key",
			values: map[string]string{"foo": "hello"},
			result: map[string]string{"foo": "hello"},
		},
		{
			title:  "bool true values are normalized",
			values: map[string]string{"egg": "yes"},
			result: map[string]string{"egg": "true"},
		},
		{
			title:  "multiple values are all loaded",
			values: map[string]string{"foo": "x", "egg": "1"},
			result: map[string]string{"foo": "x", "egg": "true"},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			m, err := New(schema, c.values)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}

			for name, value := range c.result {
				raw, err := m.GetRaw(name)
				if err != nil {
					t.Errorf("expected err to be nil: %v", err)
				}
				if expected, actual := value, raw; expected != actual {
					t.Errorf("expected: %v, actual: %v", expected, actual)
				}
			}
		})
	}
}

// If some keys fail to load, an ErrorList with the offending issues is
// returned.
func TestLoadWithError(t *testing.T) {
	var cases = []struct {
		title   string
		schema  Schema            // Test schema to use
		values  map[string]string // Initial values
		message string            // Expected error message
	}{
		{
			title:   "schema has no key with the given name",
			schema:  Schema{},
			values:  map[string]string{"bar": ""},
			message: "cannot set 'bar' to '': unknown key \"bar\"",
		},
		{
			title:   "validation fails",
			schema:  Schema{"foo": NewBoolKey()},
			values:  map[string]string{"foo": "yyy"},
			message: "cannot set 'foo' to 'yyy': invalid boolean",
		},
		{
			title:   "only the first of multiple errors is shown (in key name order)",
			schema:  Schema{"foo": NewBoolKey()},
			values:  map[string]string{"foo": "yyy", "bar": ""},
			message: "cannot set 'bar' to '': unknown key \"bar\" (and 1 more errors)",
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			_, err := New(c.schema, c.values)
			if expected, actual := c.message, err.Error(); expected != actual {
				t.Errorf("expected: %s, actual: %s", expected, actual)
			}
		})
	}
}

// Changing a config Map mutates the initial values.
func TestChange(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": {Setter: upperCase},
		"egg": NewBoolKey(),
		"yuk": NewBoolKey(WithDefault("true")),
		"xyz": {Hidden: true},
	}
	values := map[string]string{ // Initial values
		"foo": "hello",
		"bar": "x",
		"xyz": "sekret",
	}

	cases := []struct {
		title  string
		values map[string]interface{} // New values
		result map[string]string      // Expected values after change
	}{
		{
			"plain change of regular key",
			map[string]interface{}{"foo": "world"},
			map[string]string{"foo": "world"},
		},
		{
			"key setter is honored",
			map[string]interface{}{"bar": "y"},
			map[string]string{"bar": "Y"},
		},
		{
			"bool true values are normalized",
			map[string]interface{}{"egg": "yes"},
			map[string]string{"egg": "true"},
		},
		{
			"bool false values are normalized",
			map[string]interface{}{"yuk": "0"},
			map[string]string{"yuk": "false"},
		},
		{
			"the special value 'true' is a passthrough for hidden keys",
			map[string]interface{}{"xyz": true},
			map[string]string{"xyz": "sekret"},
		},
		{
			"the special value nil is converted to empty string",
			map[string]interface{}{"foo": nil},
			map[string]string{"foo": ""},
		},
		{
			"multiple values are all mutated",
			map[string]interface{}{"foo": "x", "bar": "hey", "egg": "0"},
			map[string]string{"foo": "x", "bar": "HEY", "egg": "false"},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			m, err := New(schema, values)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}

			_, err = m.Change(c.values)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}

			for name, value := range c.result {
				raw, err := m.GetRaw(name)
				if err != nil {
					t.Errorf("expected err to be nil: %v", err)
				}
				if expected, actual := value, raw; expected != actual {
					t.Errorf("expected: %s, actual: %s", expected, actual)
				}
			}
		})
	}
}

// A map of changed key/value pairs is returned.
func TestMapWithChangeReturnsChangedKeys(t *testing.T) {
	schema := Schema{
		"foo": NewBoolKey(WithDefault("false")),
		"bar": {Default: "egg"},
	}
	values := map[string]string{"foo": "true"} // Initial values

	cases := []struct {
		title   string
		changes map[string]interface{} // New values
		changed map[string]string      // Keys that should have actually changed
	}{
		{
			title:   "plain single change",
			changes: map[string]interface{}{"foo": "no"},
			changed: map[string]string{"foo": "false"},
		},
		{
			title:   "unchanged boolean value, even if it's spelled 'yes' and not 'true'",
			changes: map[string]interface{}{"foo": "yes"},
			changed: map[string]string{},
		},
		{
			title:   "unset value",
			changes: map[string]interface{}{"foo": ""},
			changed: map[string]string{"foo": "false"},
		},
		{
			title:   "unchanged value, since it matches the default",
			changes: map[string]interface{}{"foo": "true", "bar": "egg"},
			changed: map[string]string{},
		},
		{
			title:   "multiple changes",
			changes: map[string]interface{}{"foo": "false", "bar": "baz"},
			changed: map[string]string{"foo": "false", "bar": "baz"},
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			m, err := New(schema, values)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}

			changed, err := m.Change(c.changes)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}
			if expected, actual := c.changed, changed; !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

// If some keys fail to load, an ErrorList with the offending issues is
// returned.
func TestMapWithChangeError(t *testing.T) {
	schemaFoo := Schema{
		"foo": NewBoolKey(WithDefault("false")),
	}

	var cases = []struct {
		schema  Schema
		title   string
		changes map[string]interface{}
		message string
	}{
		{
			schema:  schemaFoo,
			title:   "schema has no key with the given name",
			changes: map[string]interface{}{"xxx": ""},
			message: "cannot set 'xxx' to '': unknown key \"xxx\"",
		},
		{
			schema:  schemaFoo,
			title:   "validation fails",
			changes: map[string]interface{}{"foo": "yyy"},
			message: "cannot set 'foo' to 'yyy': invalid boolean",
		},
		{
			schema: Schema{
				"egg": {Setter: failingSetter},
			},
			title:   "custom setter fails",
			changes: map[string]interface{}{"egg": "xxx"},
			message: "cannot set 'egg' to 'xxx': boom",
		},
		{
			schema: Schema{
				"egg": NewBoolKey(),
			},
			title:   "non string value",
			changes: map[string]interface{}{"egg": 123},
			message: "cannot set 'egg': invalid type int",
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			m, err := New(c.schema, nil)
			if err != nil {
				t.Errorf("expected err to be nil: %v", err)
			}

			_, err = m.Change(c.changes)
			if expected, actual := c.message, err.Error(); expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

// A Map dump contains only values that differ from their default. Hidden
// values are obfuscated.
func TestMapWithDump(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": {Default: "x"},
		"egg": {Hidden: true},
	}
	values := map[string]string{
		"foo": "hello",
		"bar": "x",
		"egg": "123",
	}
	m, err := New(schema, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	dump := map[string]interface{}{
		"foo": "hello",
	}
	got, err := m.Dump(false)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := dump, got; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestMapWithDumpIncludeDefault(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": {Default: "x"},
		"egg": {Hidden: true},
	}
	values := map[string]string{
		"foo": "hello",
		"bar": "x",
		"egg": "123",
	}
	m, err := New(schema, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	dump := map[string]interface{}{
		"bar": "x",
		"foo": "hello",
	}
	got, err := m.Dump(true)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := dump, got; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestMapWithGetString(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": NewBoolKey(),
		"egg": NewInt64Key(),
	}
	values := map[string]string{
		"foo": "hello",
		"bar": "true",
		"egg": "123",
	}

	m, err := New(schema, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	value, err := m.GetString("foo")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := "hello", value; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestMapWithGetBool(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": NewBoolKey(),
		"egg": NewInt64Key(),
	}
	values := map[string]string{
		"foo": "hello",
		"bar": "true",
		"egg": "123",
	}

	m, err := New(schema, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	value, err := m.GetBool("bar")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := true, value; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestMapWithGetInt64(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": NewBoolKey(),
		"egg": NewInt64Key(),
	}
	values := map[string]string{
		"foo": "hello",
		"bar": "true",
		"egg": "123",
	}

	m, err := New(schema, values)
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}

	value, err := m.GetInt64("egg")
	if err != nil {
		t.Errorf("expected err to be nil: %v", err)
	}
	if expected, actual := int64(123), value; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// A key setter that always fail.
func failingSetter(string) (string, error) {
	return "", errors.Errorf("boom")
}

// A key setter that uppercases the value.
func upperCase(v string) (string, error) {
	return strings.ToUpper(v), nil
}
