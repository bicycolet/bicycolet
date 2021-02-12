package provision

import (
	"fmt"
	"testing"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tickers/stub"
	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
)

func TestProvision(t *testing.T) {
	t.Parallel()

	t.Run("take", func(t *testing.T) {
		tokens := New(100, time.Millisecond, &stubTicker{})

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("take the fill", func(t *testing.T) {
		ticker := &stubTicker{}
		tokens := New(2, time.Millisecond, ticker)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(0), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}

		ticker.Ticker().Advance()

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("take with no frequency", func(t *testing.T) {
		ticker := &stubTicker{}
		tokens := New(100, -1, ticker)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})
}

func ExampleToken() {
	ticker := &stubTicker{}
	tokens := New(2, time.Millisecond, ticker)
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(1))
	ticker.Ticker().Advance()
	fmt.Println(tokens.Take(1))

	// Output:
	// 1
	// 1
	// 0
	// 1
}

type stubTicker struct {
	ticker *stub.Ticker
}

func (s *stubTicker) New(_ time.Duration, fn func()) ticker.Ticker {
	s.ticker = stub.New(fn)
	return s.ticker
}

func (s *stubTicker) Ticker() *stub.Ticker {
	return s.ticker
}
