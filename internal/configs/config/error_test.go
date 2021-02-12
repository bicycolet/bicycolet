package config

import "testing"

// Errors can be sorted by key name, and the global error message mentions the
// first of them.
func TestError(t *testing.T) {
	t.Parallel()

	t.Run("error list", func(t *testing.T) {
		var errors ErrorList
		errors.Add("foo", "xxx", "boom")
		errors.Add("bar", "yyy", "ugh")
		errors.Sort()

		wanted := "cannot set 'bar' to 'yyy': ugh (and 1 more errors)"
		if expected, actual := wanted, errors.Error(); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("no errors", func(t *testing.T) {
		var errors ErrorList
		errors.Sort()

		wanted := "no errors"
		if expected, actual := wanted, errors.Error(); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("one error", func(t *testing.T) {
		var errors ErrorList
		errors.Add("foo", "xxx", "boom")
		errors.Sort()

		wanted := "cannot set 'foo' to 'xxx': boom"
		if expected, actual := wanted, errors.Error(); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}
