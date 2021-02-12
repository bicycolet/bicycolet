package provision

import (
	"fmt"
	"testing"
	"time"
)

func TestProvision(t *testing.T) {
	t.Parallel()

	t.Run("take", func(t *testing.T) {
		tokens := New(100, time.Millisecond)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("take the fill", func(t *testing.T) {
		tokens := New(2, time.Millisecond)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
		if expected, actual := int64(0), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}

		time.Sleep(510 * time.Millisecond)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("take with no frequency", func(t *testing.T) {
		tokens := New(100, -1)

		if expected, actual := int64(1), tokens.Take(1); expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})
}

func ExampleProvision() {
	tokens := New(2, time.Millisecond)
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(1))
	fmt.Println(tokens.Take(1))
	time.Sleep(510 * time.Millisecond)
	fmt.Println(tokens.Take(1))

	// Output:
	// 1
	// 1
	// 0
	// 1
}
