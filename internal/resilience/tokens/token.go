package tokens

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tickers"
	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
	"github.com/bicycolet/bicycolet/internal/resilience/tokens/bucket"
	"github.com/bicycolet/bicycolet/internal/resilience/tokens/provision"
	"github.com/bicycolet/bicycolet/internal/resilience/tokens/token"
	"github.com/pkg/errors"
)

// TokenType defines the token we want to use for resilience.
type TokenType int

const (
	// Provision token that describes what token to use.
	Provision TokenType = iota

	// Bucket token describes a token bucket that is stable.
	Bucket
)

// New creates a new encoding token based on the type.
func New(t TokenType, capacity int64, freq time.Duration, options ...Option) (token.Token, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	switch t {
	case Provision:
		return provision.New(capacity, freq, opts.ticker), nil
	case Bucket:
		return bucket.New(capacity), nil
	default:
		return nil, errors.Errorf("invalid token type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	ticker token.Ticker
}

// WithTicker sets the ticker for token buckets that need a timer.
func WithTicker(ticker token.Ticker) Option {
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
