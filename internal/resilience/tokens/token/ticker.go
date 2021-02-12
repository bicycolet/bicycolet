package token

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
)

// Ticker defines a way to get a new ticker.
type Ticker interface {
	// New creates a new ticker using the time duration.
	New(time.Duration, func()) ticker.Ticker
}
