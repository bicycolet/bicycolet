package generic

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bicycolet/bicycolet/internal/configs/config"
	"github.com/pkg/errors"
)

// Config is a structured map of config keys to config values.
//
// Each legal key is declared in a config Schema using a key object.
type Config struct {
	schema   config.Schema
	values   map[string]string
	hydrated bool
}

// New creates a new configuration Map with the given schema and initial
// values. It is meant to be called with a set of initial values that were set
// at a previous time and persisted to some storage like a database.
//
// If one or more keys fail to be loaded, return an ErrorList describing what
// went wrong. Non-failing keys are still loaded in the returned Map.
func New(schema config.Schema, values map[string]string) (*Config, error) {
	c := &Config{
		schema: schema,
		values: make(map[string]string),
	}

	// Populate the initial values.
	_, err := c.update(values)
	return c, err
}

// Change the values of this configuration Map.
//
// Return a map of key/value pairs that were actually changed. If
// some keys fail to apply, details are included in the returned
// ErrorList.
func (c *Config) Change(changes map[string]interface{}) (map[string]string, error) {
	values := make(map[string]string, len(c.schema))

	var errs config.ErrorList
	for name, change := range changes {
		key, ok := c.schema[name]

		// When a hidden value is set to "true" in the change set, it
		// means "keep it unchanged", so we replace it with our current
		// value.
		if ok && key.Hidden && change == true {
			var err error
			if change, err = c.Key(name); err != nil {
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
		case time.Duration:
			values[name] = v.String()
		case float64:
			// We don't actually support floats yet, but because our API is JSON
			// we can end up with floats here instead of ints.
			values[name] = strconv.FormatInt(int64(v), 10)
		default:
			errs.Add(name, nil, fmt.Sprintf("invalid type %T", v))
		}
	}

	// Any key not explicitly set, is considered unset.
	for name, key := range c.schema {
		if _, ok := values[name]; !ok {
			values[name] = key.Default
		}
	}

	if errs.Len() > 0 {
		return nil, errs
	}

	names, err := c.update(values)

	changed := make(map[string]string)
	for _, name := range names {
		changed[name], err = c.Key(name)
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
func (c *Config) Dump(includeDefault bool) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	for name, key := range c.schema {
		value, err := c.Key(name)
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

// Key returns the value of the given key, which must be of type String.
func (c *Config) Key(name string) (string, error) {
	key, err := c.schema.GetKey(name)
	if err != nil {
		return "", errors.WithStack(err)
	}
	value, ok := c.values[name]
	if !ok {
		value = key.Default
	}
	return value, nil
}

// Accessor returns the accessor for a given config.
func (c *Config) Accessor() *config.Accessor {
	return config.NewAccessor(c, c.schema)
}

// Update the current values in the map using the newly provided ones. Return a
// list of key names that were actually changed and an ErrorList with possible
// errors.
func (c *Config) update(values map[string]string) ([]string, error) {
	defer func() {
		c.hydrated = true
	}()
	// Update our keys with the values from the given map, and keep track
	// of which keys actually changed their value.
	var (
		errs  config.ErrorList
		names []string
	)
	for name, value := range values {
		changed, err := c.set(name, value, !c.hydrated)
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
func (c *Config) set(name string, value string, initial bool) (bool, error) {
	key, ok := c.schema[name]
	if !ok {
		return false, errors.Errorf("unknown key %q", name)
	}

	if err := key.Validate(value); err != nil {
		return false, err
	}

	current, err := c.Key(name)
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
		delete(c.values, name)
	} else {
		c.values[name] = value
	}

	return true, nil
}
