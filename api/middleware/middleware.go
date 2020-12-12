package middleware

import (
	"net/http"
	"strings"

	"github.com/bicycolet/bicycolet/internal/clock"
	"github.com/go-kit/kit/log"
)

// Func defines the middleware intended to run durring a http request.
type Func func(http.Handler, log.Logger) http.Handler

// Builder builds all the middleware.
type Builder struct {
	logger      log.Logger
	middlewares []Func
}

// New creates a middleware runner
func New(logger log.Logger) *Builder {
	return &Builder{
		logger: logger,
	}
}

// Add a new middleware to be executed.
func (m *Builder) Add(fn Func) {
	m.middlewares = append(m.middlewares, fn)
}

// Build creates one middleware from a series of middlewares
func (m *Builder) Build() http.Handler {
	// Ensure we always have a middleware to run.
	var last http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Considering this is the last handler, we can actually time the LOT!
	})
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		last = m.middlewares[i](last, m.logger)
	}
	return last
}

// InternalHandler allows the handling of internal endpoint handling.
type InternalHandler struct {
	next http.Handler
	fn   http.HandlerFunc
}

// InternalHandlerFunc handles any internal routing to prevent additional logic
// to handle internal error checking.
func InternalHandlerFunc(next http.Handler, fn http.HandlerFunc) InternalHandler {
	return InternalHandler{
		next: next,
		fn:   fn,
	}
}

// ServeHTTP calls f(w, r).
func (h InternalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/internal") {
		h.fn(w, r)
		return
	}

	h.next.ServeHTTP(w, r)
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	logger log.Logger
	clock  clock.Clock
}

// WithLogger sets the logger on the option
func WithLogger(logger log.Logger) Option {
	return func(options *options) {
		options.logger = logger
	}
}

// WithClock sets the clock on the options
func WithClock(clock clock.Clock) Option {
	return func(options *options) {
		options.clock = clock
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		logger: log.NewNopLogger(),
		clock:  clock.WallClock{},
	}
}
