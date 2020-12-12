package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tsenart/tb"
)

// Throttle is a middleware to throttle the
func Throttle(rate int64, freq time.Duration) Func {
	throttler := tb.NewThrottler(freq)
	return func(next http.Handler, logger log.Logger) http.Handler {
		return InternalHandlerFunc(next, func(w http.ResponseWriter, r *http.Request) {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			if throttler.Halt(host, 1, rate) {
				http.Error(w, "Too many requests", 429)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
