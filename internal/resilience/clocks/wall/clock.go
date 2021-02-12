package wall

import (
	"sync/atomic"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/clocks/clock"
)

// Clock represents a wall clock.
type Clock struct {
	stamp uint64
}

// New creates a new time clock
func New() *Clock {
	return &Clock{
		stamp: uint64(time.Now().UnixNano()),
	}
}

// Now is used to return the current value of the clock
func (l *Clock) Now() clock.Time {
	return Time(atomic.LoadUint64(&l.stamp))
}

// Increment is used to increment and return the value of the clock
func (l *Clock) Increment() clock.Time {
	return Time(atomic.AddUint64(&l.stamp, 1))
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *Clock) Witness(v clock.Time) {
WITNESS:
	// if the other value is old, we do not need to do anything
	var (
		cur   = atomic.LoadUint64(&l.stamp)
		other = v.Value()
	)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.stamp, cur, other+1) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
}

// Clone a clock with the same local time underneath
func (l *Clock) Clone() clock.Clock {
	return &Clock{stamp: atomic.LoadUint64(&l.stamp)}
}

// Reset the clock
func (l *Clock) Reset() {
RESET:
	var (
		cur   = atomic.LoadUint64(&l.stamp)
		other = uint64(time.Now().UnixNano())
	)
	if other == cur {
		return
	}

	if !atomic.CompareAndSwapUint64(&l.stamp, cur, other) {
		goto RESET
	}
}

// Time defines a unit of time within a lamport clock.
type Time uint64

// Value represents the underling unit within a clock.
func (t Time) Value() uint64 {
	return uint64(t)
}

// Before returns if the current unit of time is before the other.
func (t Time) Before(other clock.Time) bool {
	return uint64(t) < uint64(other.Value())
}

// After returns if the current unit of time is after the other.
func (t Time) After(other clock.Time) bool {
	return uint64(t) > uint64(other.Value())
}
