package config

import (
	"time"

	"github.com/pkg/errors"
)

// Accessor retrieves values.
type Accessor struct {
	config Config
	schema Schema
}

// NewAccessor creates a new accessor to get values out of the config in a typed
// way.
func NewAccessor(c Config, s Schema) *Accessor {
	return &Accessor{
		config: c,
		schema: s,
	}
}

// String returns the value of the given key, which must be of type String.
func (a *Accessor) String(name string) (string, error) {
	v, err := a.applyGetter(name)
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

// Bool returns the value of the given key, which must be of type Bool.
func (a *Accessor) Bool(name string) (bool, error) {
	v, err := a.applyGetter(name)
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

// Int64 returns the value of the given key, which must be of type Int64.
func (a *Accessor) Int64(name string) (int64, error) {
	v, err := a.applyGetter(name)
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

// Duration returns the value of the given key, which must be of type Duration.
func (a *Accessor) Duration(name string) (time.Duration, error) {
	v, err := a.applyGetter(name)
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

func (a *Accessor) applyGetter(name string) (interface{}, error) {
	key, err := a.schema.GetKey(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	value, err := a.config.Key(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if key.Getter != nil {
		return key.Getter(value)
	}

	return interface{}(value), nil
}
