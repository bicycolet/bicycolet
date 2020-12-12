package middleware

import (
	"errors"
	"net/http"

	"github.com/bicycolet/bicycolet/api/responses"
	"github.com/go-kit/kit/log"
)

// Setup defines an interface for waiting for a setup channel to be closed.
type Setup interface {
	SetupChan() <-chan struct{}
}

// WaitReady is a middleware to check if the daemon is setup before letting
// any routes through.
func WaitReady(setup Setup) Func {
	return func(next http.Handler, logger log.Logger) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Block public API requests until we're done with basic
			// initialization tasks.
			select {
			case <-setup.SetupChan():
				next.ServeHTTP(w, r)
			default:
				responses.Unavailable(errors.New("waiting for setup to finish")).Render(logger, w)
				return
			}
		})
	}
}
