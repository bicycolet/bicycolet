package circuit

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/breakers/breaker"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/stub"
	"github.com/bicycolet/bicycolet/internal/resilience/timers/timer"
)

func TestBreaker(t *testing.T) {
	t.Parallel()

	t.Run("func called", func(t *testing.T) {
		called := false

		circuit := New(1, time.Second, &stubTimer{})
		err := circuit.Run(func() error {
			called = true
			return nil
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("func returns error", func(t *testing.T) {
		called := false

		circuit := New(1, time.Second, &stubTimer{})
		err := circuit.Run(func() error {
			called = true
			return errors.New("bad")
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := false, breaker.ErrBreakerOpen(err); expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("second call isn't successful", func(t *testing.T) {
		called := false

		circuit := New(1, time.Second, &stubTimer{})
		_ = circuit.Run(func() error {
			return errors.New("bad")
		})

		err := circuit.Run(func() error {
			called = true
			return nil
		})

		if expected, actual := false, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, breaker.ErrBreakerOpen(err); expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})

	t.Run("second call is successful after timeout", func(t *testing.T) {
		called := false

		expiryTimeout := 100 * time.Millisecond

		timer := &stubTimer{}
		circuit := New(1, expiryTimeout, timer)
		_ = circuit.Run(func() error {
			return errors.New("bad")
		})

		timer.Timer().Advance(expiryTimeout + (time.Millisecond * 4))

		err := circuit.Run(func() error {
			called = true
			return nil
		})

		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %t, actual: %t, err: %v", expected, actual, err)
		}
	})
}

func TestBreakerTransitions(t *testing.T) {
	t.Parallel()

	timer := &stubTimer{}
	circuit := New(3, 100*time.Millisecond, timer)

	// Make sure we trigger an open circuit
	for i := 0; i < 3; i++ {
		_ = circuit.Run(func() error {
			return errors.New("bad")
		})
	}

	// Make sure that we're still open
	if err := circuit.Run(func() error {
		return nil
	}); err == nil {
		t.Fatal("circuit breaker should be open, but isn't!")
	}

	timer.Timer().Advance(110 * time.Millisecond)

	if err := circuit.Run(func() error {
		return nil
	}); err != nil {
		t.Fatal("circuit breaker should be closed, but isn't!")
	}

	// Make sure we trigger an open circuit
	for i := 0; i < 3; i++ {
		_ = circuit.Run(func() error {
			return errors.New("bad")
		})
	}

	timer.Timer().Advance(110 * time.Millisecond)

	if err := circuit.Run(func() error {
		return errors.New("bad")
	}); err == nil {
		t.Fatal("circuit breaker should be open, but isn't!")
	}

	timer.Timer().Advance(110 * time.Millisecond)

	if err := circuit.Run(func() error {
		return nil
	}); err != nil {
		t.Fatal("circuit breaker should be closed, but isn't!")
	}
}

func Example() {
	circuit := New(10, 10*time.Millisecond, &stubTimer{})
	err := circuit.Run(func() error {
		// communicate with an external service
		return nil
	})

	switch {
	case err == nil:
		fmt.Println("success!")
	case breaker.ErrBreakerOpen(err):
		fmt.Println("breaker is open")
	default:
		fmt.Println("other error")
	}

	// Output:
	// success!
}

type stubTimer struct {
	timer *stub.Timer
}

func (s *stubTimer) New(expiry time.Duration, fn func()) timer.Timer {
	s.timer = stub.New(expiry, fn)
	return s.timer
}

func (s *stubTimer) Timer() *stub.Timer {
	return s.timer
}
