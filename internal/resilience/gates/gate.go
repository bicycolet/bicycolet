package gates

import (
	"github.com/bicycolet/bicycolet/internal/resilience/gates/gate"
	"github.com/pkg/errors"
)

// GateType defines the ticker we want to use for resilience.
type GateType int

const (
	// Branch ticker that describes what ticker to use.
	Branch GateType = iota
)

// New creates a new encoding ticker based on the type.
func New(t GateType, left, right func() error) (gate.Gate, error) {
	switch t {
	case Branch:
		return branch.New(left, right), nil
	default:
		return nil, errors.Errorf("invalid ticker type %q", t)
	}
}
