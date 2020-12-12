package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestSchemaWithDefaults(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": {Default: "x"},
	}
	values := map[string]interface{}{"foo": "", "bar": "x"}
	if expected, actual := values, schema.Defaults(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestSchemaWithKeys(t *testing.T) {
	schema := Schema{
		"foo": {},
		"bar": {Default: "x"},
	}
	keys := []string{"bar", "foo"}
	if expected, actual := keys, schema.Keys(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// Exercise valid values.
func TestKeyWithValidate(t *testing.T) {
	for k, c := range []struct {
		node  key
		value string
	}{
		{node: key{}, value: "hello"},
		{node: NewBoolKey(), value: "yes"},
		{node: NewBoolKey(), value: "0"},
		{node: NewInt64Key(), value: "666"},
		{node: NewBoolKey(), value: ""},
		{node: NewStringKey(WithDefault("foo"), WithValidator(isNotEmptyString)), value: ""},
	} {
		t.Run(fmt.Sprintf("validate %d", k), func(t *testing.T) {
			if err := c.node.validate(c.value); err != nil {
				t.Errorf("expected err to be nil: got %v", err)
			}
		})
	}
}

// Validator that returns an error if the value is not the empty string.
func isNotEmptyString(value string) error {
	if value == "" {
		return errors.Errorf("empty value not valid")
	}
	return nil
}

// Exercise all possible validation errors.
func TestKey_validateError(t *testing.T) {
	for _, c := range []struct {
		node    key
		value   string
		message string
	}{
		{node: NewInt64Key(), value: "1.2", message: "invalid integer"},
		{node: NewBoolKey(), value: "yyy", message: "invalid boolean"},
		{node: key{Validator: func(string) error { return errors.Errorf("ugh") }}, value: "", message: "ugh"},
		{node: key{Deprecated: "don't use this"}, value: "foo", message: "deprecated: don't use this"},
	} {
		t.Run(c.message, func(t *testing.T) {
			if err := c.node.validate(c.value); err == nil {
				t.Errorf("expected err to not be nil")
			}
		})
	}
}
