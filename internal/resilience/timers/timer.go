package timers

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/wall"
	"github.com/pkg/errors"
)

// TimerType defines the timer we want to use for resilience.
type TimerType int

const (
	// Wall timer that describes what timer to use.
	Wall TimerType = iota
)

// New creates a new encoding timer based on the type.
func New(t TimerType, expiry time.Duration, fn func()) (timer.Timer, error) {
	switch t {
	case Wall:
		return wall.New(expiry, fn), nil
	default:
		return nil, errors.Errorf("invalid timer type %q", t)
	}
}
