package branch

import "sync/atomic"

const (
	leftState uint32 = iota
	rightState
)

// Gate defines a way to switch branches before running.
type Gate struct {
	left, right func() error
	state       uint32
}

// New creates a new gate.
func New(left, right func() error) *Gate {
	return &Gate{
		left:  left,
		right: right,
	}
}

// Switch between gates.
func (g *Gate) Switch() bool {
SWITCH:
	var (
		state = atomic.LoadUint32(&g.state)
		next  uint32
	)
	switch state {
	case leftState:
		next = rightState
	case rightState:
		next = leftState
	}

	// Ensure that our gate has actually moved.
	if !atomic.CompareAndSwapUint32(&g.state, state, next) {
		// The CAS failed, so we just retry. Eventually our CAS should
		// succeed or a future switch will pass.
		goto SWITCH
	}
	return next == leftState
}

// Run the gate.
func (g *Gate) Run() error {
	state := atomic.LoadUint32(&g.state)
	if state == leftState {
		return g.left()
	}
	return g.right()
}
