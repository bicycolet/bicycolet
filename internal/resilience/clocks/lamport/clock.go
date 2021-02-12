package lamport

import (
	"sync/atomic"

	"github.com/bicycolet/bicycolet/internal/resilience/clocks/clock"
)

// Clock represents a lamport clock.
type Clock struct {
	counter uint64
}

// New creates a new lamport clock
func New() *Clock {
	return &Clock{
		counter: 0,
	}
}

// Now is used to return the current value of the clock
func (l *Clock) Now() clock.Time {
	return Time(atomic.LoadUint64(&l.counter))
}

// Increment is used to increment and return the value of the clock
func (l *Clock) Increment() clock.Time {
	return Time(atomic.AddUint64(&l.counter, 1))
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *Clock) Witness(v clock.Time) {
WITNESS:
	// If the other value is old, we do not need to do anything
	var (
		cur   = atomic.LoadUint64(&l.counter)
		other = v.Value()
	)
	if other < cur {
		return
	}

	// Ensure that our local clock is at least one ahead.
	if !atomic.CompareAndSwapUint64(&l.counter, cur, other+1) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future witness will pass us by and our witness
		// will end.
		goto WITNESS
	}
}

// Clone a clock with the same local time underneath
func (l *Clock) Clone() clock.Clock {
	return &Clock{counter: atomic.LoadUint64(&l.counter)}
}

// Reset the clock
func (l *Clock) Reset() {
RESET:
	cur := atomic.LoadUint64(&l.counter)
	if cur == 0 {
		return
	}

	if !atomic.CompareAndSwapUint64(&l.counter, cur, 0) {
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
