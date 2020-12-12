package middleware

import (
	"net/http"

	"github.com/go-kit/kit/log"
)

// ContentType is a middleware to check if the daemon is setup before letting
// any routes through.
func ContentType() Func {
	return func(next http.Handler, logger log.Logger) http.Handler {
		return InternalHandlerFunc(next, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
