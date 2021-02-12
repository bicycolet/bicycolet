package tickers

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
	"github.com/bicycolet/bicycolet/internal/resilience/tickers/wall"
	"github.com/pkg/errors"
)

// TickerType defines the ticker we want to use for resilience.
type TickerType int

const (
	// Wall ticker that describes what ticker to use.
	Wall TickerType = iota
)

// New creates a new encoding ticker based on the type.
func New(t TickerType, expiry time.Duration, fn func()) (ticker.Ticker, error) {
	switch t {
	case Wall:
		return wall.New(expiry, fn), nil
	default:
		return nil, errors.Errorf("invalid ticker type %q", t)
	}
}
