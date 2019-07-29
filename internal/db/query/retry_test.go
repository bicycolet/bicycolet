package query_test

import (
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/query"
	"github.com/pkg/errors"
)

func TestIsRetriableError(t *testing.T) {
	ok := query.IsRetriableError(errors.New("bad"))
	if ok {
		t.Errorf("expected ok to be false")
	}
}

func TestIsRetriableErrorWithNil(t *testing.T) {
	ok := query.IsRetriableError(nil)
	if ok {
		t.Errorf("expected ok to be false")
	}
}
