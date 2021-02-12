package configs

import (
	"github.com/bicycolet/bicycolet/internal/configs/config"
	"github.com/bicycolet/bicycolet/internal/configs/generic"
	"github.com/pkg/errors"
)

// ConfigType defines the token we want to use for resilience.
type ConfigType int

const (
	// Map config that describes what config backing to use.
	Map ConfigType = iota
)

// New creates a new encoding token based on the type.
func New(t ConfigType, schema config.Schema, values map[string]string) (config.Config, error) {
	switch t {
	case Map:
		return generic.New(schema, values)
	default:
		return nil, errors.Errorf("invalid config type %q", t)
	}
}
