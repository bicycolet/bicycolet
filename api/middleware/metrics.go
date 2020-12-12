package middleware

import (
	"context"
	"net/http"

	"github.com/bicycolet/bicycolet/internal/instrumentation"
	"github.com/go-kit/kit/log"
)

// Instruments defines the instruments to be injected into the stack.
type Instruments interface {
	Gauge(string) instrumentation.Gauge
	Summary(string) instrumentation.Summary
	SummaryVec(string, ...string) instrumentation.SummaryVec
}

// Metrics is a middleware to check if the daemon is setup before letting
// any routes through.
func Metrics(instruments Instruments) Func {
	return func(next http.Handler, logger log.Logger) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, instrumentation.Metrics, instruments)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
