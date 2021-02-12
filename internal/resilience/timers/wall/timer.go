package wall

import (
	"sync"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
)

// Timer represents a wall timer.
type Timer struct {
	mutex   sync.RWMutex
	current time.Time
	expiry  time.Duration
	timer   *time.Timer
}

// New creates a wall timer.
func New(expiry time.Duration, fn func()) *Timer {
	return &Timer{
		current: time.Now(),
		expiry:  expiry,
		timer:   time.AfterFunc(expiry, fn),
	}
}

// Now is used to return the current value of the timer
func (t *Timer) Now() timer.Time {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.current = time.Now()
	return Time(uint64(t.current.UnixNano()))
}

// After is used to know if the time now is after timer has ticked
func (t *Timer) After() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	expiry := t.current.Add(t.expiry)
	return time.Now().After(expiry)
}

// Reset the timer
func (t *Timer) Reset() error {
	t.timer.Stop()
	t.timer.Reset(t.expiry)
	return nil
}

// Time is the value of a Timer.
type Time uint64

// Value returns the underlying unit of time.
func (t Time) Value() uint64 {
	return uint64(t)
}
