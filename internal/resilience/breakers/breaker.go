package breakers

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/breakers/breaker"
	"github.com/bicycolet/bicycolet/internal/resilience/breakers/circuit"
	"github.com/bicycolet/bicycolet/internal/resilience/timers"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
	"github.com/pkg/errors"
)

// BreakerType defines the circuit we want to use for resilience.
type BreakerType int

const (
	// Circuit circuit that describes what circuit to use.
	Circuit BreakerType = iota
)

// New creates a new encoding circuit based on the type.
func New(t BreakerType, failures uint64, expiry time.Duration, options ...Option) (breaker.Breaker, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	switch t {
	case Circuit:
		return circuit.New(failures, expiry, opts.timer), nil
	default:
		return nil, errors.Errorf("invalid ticker type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	timer breaker.Timer
}

// WithTimer sets the timer for token buckets that need a timer.
func WithTimer(timer breaker.Timer) Option {
	return func(options *options) {
		options.timer = timer
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		timer: defaultTimer{},
	}
}

type defaultTimer struct{}

func (defaultTimer) New(d time.Duration, fn func()) timer.Timer {
	t, err := timers.New(timers.Wall, d, fn)
	if err != nil {
		panic("programming error: expected timer.")
	}
	return t
}
