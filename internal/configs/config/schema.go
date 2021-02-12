package config

import (
	"sort"

	"github.com/pkg/errors"
)

// Schema defines the available keys of a config Map, along with the types
// and options for their values, expressed using key objects.
type Schema map[string]Key

// Keys returns all keys defined in the schema
func (s Schema) Keys() []string {
	var i int
	keys := make([]string, len(s))
	for key := range s {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// Defaults returns a map of all key names in the schema along with their default
// values.
func (s Schema) Defaults() map[string]interface{} {
	values := make(map[string]interface{}, len(s))
	for name, key := range s {
		values[name] = key.Default
	}
	return values
}

// GetKey retrives the key associated with the given name.
func (s Schema) GetKey(name string) (Key, error) {
	k, ok := s[name]
	if !ok {
		return Key{}, errors.Errorf("attempt to access unknown key %q", name)
	}
	return k, nil
}
