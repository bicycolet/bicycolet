package stub

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
)

// Timer represents a wall timer.
type Timer struct {
	fn      func()
	expiry  time.Duration
	current time.Time
}

// New creates a wall timer.
func New(expiry time.Duration, fn func()) *Timer {
	return &Timer{
		expiry: expiry,
		fn:     fn,
	}
}

// Now is used to return the current value of the timer
func (t *Timer) Now() timer.Time {
	return Time(uint64(t.current.UnixNano()))
}

// After is used to know if the time now is after timer has ticked
func (t *Timer) After() bool {
	return t.current.UnixNano() > int64(t.expiry)
}

// Reset the timer
func (t *Timer) Reset() error {
	return nil
}

// Advance advances the time to call the next tick.
func (t *Timer) Advance(c time.Duration) {
	t.current.Add(c)
	t.fn()
}

// Time is the value of a Timer.
type Time uint64

// Value returns the underlying unit of time.
func (t Time) Value() uint64 {
	return uint64(t)
}
