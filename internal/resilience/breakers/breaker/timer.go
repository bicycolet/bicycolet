package breaker

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
)

// Timer defines a way to get a new timer.
type Timer interface {
	// New creates a new timer using the time duration.
	New(time.Duration, func()) timer.Timer
}
