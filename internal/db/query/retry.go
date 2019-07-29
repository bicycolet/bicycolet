package query

import (
	"strings"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/clock"
	"github.com/bicycolet/bicycolet/internal/resilience/retrier"
	"github.com/pkg/errors"
)

// Retry wraps a function that interacts with the database, and retries it in
// case a transient error is hit.
//
// This should by typically used to wrap transactions.
func Retry(sleeper clock.Sleeper, f func() error) error {
	retry := retrier.New(sleeper, 10, 250*time.Millisecond)
	err := retry.Run(func() error {
		err := f()
		if IsRetriableError(err) {
			return nil
		}
		return errors.WithStack(err)
	})
	return errors.WithStack(err)
}

// IsRetriableError returns true if the given error might be transient and the
// interaction can be safely retried.
func IsRetriableError(err error) bool {
	err = errors.Cause(err)
	if err == nil {
		return false
	}

	if strings.Contains(errors.Cause(err).Error(), "bad connection") {
		return true
	}
	return false
}
