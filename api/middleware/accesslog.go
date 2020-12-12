package middleware

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// AccessLog is a middleware to check if the daemon is setup before letting
// any routes through.
func AccessLog() Func {
	return func(next http.Handler, logger log.Logger) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			level.Debug(logger).Log("url", r.URL.String())
			next.ServeHTTP(w, r)
		})
	}
}
