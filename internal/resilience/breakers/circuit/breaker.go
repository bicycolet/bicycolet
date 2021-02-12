package circuit

import (
	"sync/atomic"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/breakers/breaker"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
)

const (
	openState uint32 = iota
	halfOpenState
	closedState
)

// Circuit creates a new circuit breaker.
type Circuit struct {
	state      *state
	thresholds *thresholds
	timer      timer.Timer
}

// New creates a new circuit breaker.
func New(failure uint64, expiry time.Duration, timer breaker.Timer) *Circuit {
	c := &Circuit{
		state: &state{
			value:   closedState,
			success: &guage{},
			failed:  &guage{},
		},
		thresholds: &thresholds{
			success: 1,
			failed:  failure,
		},
	}
	c.timer = timer.New(expiry, c.tick)
	return c
}

// Run a function against a given breaker.
func (c *Circuit) Run(fn func() error) error {
	state := atomic.LoadUint32(&c.state.value)
	if state == openState {
		return breaker.NewErrBreakerOpen()
	}

	err := fn()

	if c.state.failed.Current() > 0 {
		if c.timer.After() {
			c.state.Reset()
		}
	}

	switch state {
	case closedState:
		if err == nil {
			break
		}

		val := c.state.failed.Increment()
		if c.thresholds.AfterFailed(val) {
			c.open()
		} else {
			c.timer.Now()
		}

	case halfOpenState:
		if err != nil {
			c.open()
			break
		}

		val := c.state.success.Increment()
		if c.thresholds.AfterSuccess(val) {
			c.close()
		}
	}

	return err
}

func (c *Circuit) open() {
	c.state.Reset()
	atomic.StoreUint32(&c.state.value, openState)

	c.timer.Reset()
}

func (c *Circuit) close() {
	c.state.Reset()
	atomic.StoreUint32(&c.state.value, closedState)
}

func (c *Circuit) tick() {
	if state := atomic.LoadUint32(&c.state.value); state == closedState {
		return
	}

	c.state.Reset()
	atomic.StoreUint32(&c.state.value, halfOpenState)
}

type guage struct {
	counter uint64
}

func (g *guage) Current() uint64 {
	return atomic.LoadUint64(&g.counter)
}

// increment the guage.
func (g *guage) Increment() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

// Reset the guage
func (g *guage) Reset() {
RESET:
	cur := atomic.LoadUint64(&g.counter)
	if cur == 0 {
		return
	}

	if !atomic.CompareAndSwapUint64(&g.counter, cur, 0) {
		goto RESET
	}
}

type state struct {
	value           uint32
	success, failed *guage
}

func (m state) Reset() {
	m.success.Reset()
	m.failed.Reset()
}

type thresholds struct {
	success, failed uint64
}

func (t *thresholds) AfterSuccess(x uint64) bool {
	return x >= t.success
}

func (t *thresholds) AfterFailed(x uint64) bool {
	return x >= t.failed
}
