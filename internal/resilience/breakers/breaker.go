package breakers

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/breakers/breaker"
	"github.com/bicycolet/bicycolet/internal/resilience/tickers"
	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
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
		return circuit.New(failures, expiry, opts.ticker), nil
	default:
		return nil, errors.Errorf("invalid ticker type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	ticker breaker.Ticker
}

// WithTicker sets the ticker for token buckets that need a timer.
func WithTicker(ticker breaker.Ticker) Option {
	return func(options *options) {
		options.ticker = ticker
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		ticker: defaultTicker{},
	}
}

type defaultTicker struct{}

func (defaultTicker) New(d time.Duration, fn func()) ticker.Ticker {
	t, err := tickers.New(tickers.Wall, d, fn)
	if err != nil {
		panic("programming error: expected ticker.")
	}
	return t
}
