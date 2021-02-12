package clocks

import (
	"github.com/bicycolet/bicycolet/internal/resilience/clocks/clock"
	"github.com/bicycolet/bicycolet/internal/resilience/clocks/lamport"
	"github.com/bicycolet/bicycolet/internal/resilience/clocks/wall"
	"github.com/pkg/errors"
)

// ClockType defines the clock we want to use for resilience.
type ClockType int

const (
	// Lamport clock that describes what clock to use.
	Lamport ClockType = iota
	// Wall clock that describes what clock to use.
	Wall
)

// New creates a new encoding clock based on the type.
func New(t ClockType) (clock.Clock, error) {
	switch t {
	case Lamport:
		return lamport.New(), nil
	case Wall:
		return wall.New(), nil
	default:
		return nil, errors.Errorf("invalid clock type %q", t)
	}
}
