package config

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Map is a structured map of config keys to config values.
//
// Each legal key is declared in a config Schema using a key object.
type Map struct {
	schema    Schema
	values    map[string]string
	populated bool
}

// New creates a new configuration Map with the given schema and initial
// values. It is meant to be called with a set of initial values that were set
// at a previous time and persisted to some storage like a database.
//
// If one or more keys fail to be loaded, return an ErrorList describing what
// went wrong. Non-failing keys are still loaded in the returned Map.
func New(schema Schema, values map[string]string) (Map, error) {
	m := Map{
		schema: schema,
		values: make(map[string]string),
	}

	// Populate the initial values.
	_, err := m.update(values)
	return m, err
}

// Change the values of this configuration Map.
//
// Return a map of key/value pairs that were actually changed. If
// some keys fail to apply, details are included in the returned
// ErrorList.
func (m *Map) Change(changes map[string]interface{}) (map[string]string, error) {
	values := make(map[string]string, len(m.schema))

	var errs ErrorList
	for name, change := range changes {
		key, ok := m.schema[name]

		// When a hidden value is set to "true" in the change set, it
		// means "keep it unchanged", so we replace it with our current
		// value.
		if ok && key.Hidden && change == true {
			var err error
			if change, err = m.GetRaw(name); err != nil {
				errs.Add(name, nil, err.Error())
				continue
			}
		}

		// A nil object means the empty string.
		if change == nil {
			change = ""
		}

		// Sanity check that we were actually passed a string.
		switch v := change.(type) {
		case string:
			values[name] = v
		case int64:
			values[name] = strconv.FormatInt(v, 10)
		case float64:
			// We don't actually support floats yet, but because our API is JSON
			// we can end up with floats here instead of ints.
			values[name] = strconv.FormatInt(int64(v), 10)
		default:
			errs.Add(name, nil, fmt.Sprintf("invalid type %T", v))
		}
	}

	// Any key not explicitly set, is considered unset.
	for name, key := range m.schema {
		if _, ok := values[name]; !ok {
			values[name] = key.Default
		}
	}

	if errs.Len() > 0 {
		return nil, errs
	}

	names, err := m.update(values)

	changed := make(map[string]string)
	for _, name := range names {
		changed[name], err = m.GetRaw(name)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return changed, errors.WithStack(err)
}

// Dump the current configuration held by this Map.
//
// Keys that match their default value will not be included in the dump. Also,
// if a key has its Hidden attribute set to true, it will be rendered as
// "true", for obfuscating the actual value.
func (m *Map) Dump(includeDefault bool) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	for name, key := range m.schema {
		value, err := m.GetRaw(name)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if key.Hidden {
			continue
		}
		if !includeDefault && value == key.Default {
			continue
		}
		values[name] = value
	}
	return values, nil
}

// GetRaw returns the value of the given key, which must be of type String.
func (m *Map) GetRaw(name string) (string, error) {
	key, err := m.schema.getKey(name)
	if err != nil {
		return "", errors.WithStack(err)
	}
	value, ok := m.values[name]
	if !ok {
		value = key.Default
	}
	return value, nil
}

func (m *Map) applyGetter(name string) (interface{}, error) {
	key, err := m.schema.getKey(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	value, err := m.GetRaw(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if key.Getter != nil {
		return key.Getter(value)
	}

	return interface{}(value), nil
}

// GetByteArray returns the value of the given key, which must be of type String.
func (m *Map) GetByteArray(name string) ([]byte, error) {
	v, err := m.applyGetter(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch t := v.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

// GetString returns the value of the given key, which must be of type String.
func (m *Map) GetString(name string) (string, error) {
	v, err := m.applyGetter(name)
	if err != nil {
		return "", errors.WithStack(err)
	}

	switch t := v.(type) {
	case string:
		return t, nil
	default:
		return "", errors.New("unsupported type")
	}
}

// GetBool returns the value of the given key, which must be of type Bool.
func (m *Map) GetBool(name string) (bool, error) {
	v, err := m.applyGetter(name)
	if err != nil {
		return false, errors.WithStack(err)
	}

	switch t := v.(type) {
	case bool:
		return t, nil
	default:
		return false, errors.New("unsupported type")
	}
}

// GetInt64 returns the value of the given key, which must be of type Int64.
func (m *Map) GetInt64(name string) (int64, error) {
	v, err := m.applyGetter(name)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	switch t := v.(type) {
	case int64:
		return t, nil
	default:
		return -1, errors.New("unsupported type")
	}
}

// GetDuration returns the value of the given key, which must be of type Duration.
func (m *Map) GetDuration(name string) (time.Duration, error) {
	v, err := m.applyGetter(name)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	switch t := v.(type) {
	case time.Duration:
		return t, nil
	default:
		return -1, errors.New("unsupported type")
	}
}

// Clone the existing map
func (m *Map) Clone() Map {
	cloned := make(map[string]string, len(m.values))
	for k, v := range m.values {
		cloned[k] = v
	}
	return Map{
		populated: m.populated,
		schema:    m.schema,
		values:    cloned,
	}
}

// Values returns the underlying values.
func (m *Map) Values() map[string]string {
	return m.values
}

// Update the current values in the map using the newly provided ones. Return a
// list of key names that were actually changed and an ErrorList with possible
// errors.
func (m *Map) update(values map[string]string) ([]string, error) {
	defer func() {
		m.populated = true
	}()
	// Update our keys with the values from the given map, and keep track
	// of which keys actually changed their value.
	var (
		errs  ErrorList
		names []string
	)
	for name, value := range values {
		changed, err := m.set(name, value, !m.populated)
		if err != nil {
			errs.Add(name, value, err.Error())
			continue
		}
		if changed {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	var err error
	if errs.Len() > 0 {
		errs.Sort()
		err = errs
	}

	return names, err
}

// Set or change an individual key. Empty string means delete this value and
// effectively revert it to the default. Return a boolean indicating whether
// the value has changed, and error if something went wrong.
func (m *Map) set(name string, value string, initial bool) (bool, error) {
	key, ok := m.schema[name]
	if !ok {
		return false, errors.Errorf("unknown key %q", name)
	}

	err := key.validate(value)
	if err != nil {
		return false, err
	}

	current, err := m.GetRaw(name)
	if err != nil {
		return false, err
	}

	// Trigger the Setter if this key's schema has declared it.
	// Do this before change tracking takes place
	if value != "" && key.Setter != nil {
		value, err = key.Setter(value)
		if err != nil {
			return false, err
		}
	}

	// Compare the new value with the current one, and return now if they
	// are equal.
	if value == current {
		return false, nil
	}

	if value == "" {
		delete(m.values, name)
	} else {
		m.values[name] = value
	}

	return true, nil
}
