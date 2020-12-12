package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bicycolet/bicycolet/api/middleware/server"
	"github.com/bicycolet/bicycolet/api/responses"
	"github.com/bicycolet/bicycolet/internal/clock"
	"github.com/bicycolet/bicycolet/internal/instrumentation"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// ServiceRouter allows the adding of a new routable endpoint.
type ServiceRouter = server.Router

// Service represents a service that can be attached to a router.
type Service interface {
	// Subscribe to a router.
	Subscribe(ServiceRouter)
}

// Services holds the services for the router to serve
type Services struct {
	Public       []Service
	Internal     []Service
	HandlerFuncs map[string]http.HandlerFunc
}

// Router is the final middleware for routing traffic.
func Router(services Services, contextBuilder ContextBuilder, options ...Option) Func {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	return func(next http.Handler, logger log.Logger) http.Handler {
		router := mux.NewRouter()
		router.StrictSlash(false)

		router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			fmt.Fprintln(w, "OK")
		})

		for _, service := range services.Public {
			service.Subscribe(builder{
				ContextBuilder: contextBuilder,
				Prefix:         "/1.0",
				Router:         router,
				logger:         log.WithPrefix(logger, "prefix", "1.0"),
				clock:          opts.clock,
			})
		}

		for endpoint, fn := range services.HandlerFuncs {
			level.Debug(logger).Log("msg", "adding route", "pattern", endpoint)

			router.HandleFunc(endpoint, fn)
		}

		for _, service := range services.Internal {
			service.Subscribe(builder{
				ContextBuilder: contextBuilder,
				Prefix:         "/internal",
				Router:         router,
				logger:         log.WithPrefix(logger, "prefix", "internal"),
				clock:          opts.clock,
			})
		}

		router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			level.Info(logger).Log("msg", "Sending top level 404", "url", r.URL)
			w.WriteHeader(404)
			w.Header().Set("Content-Type", "application/json")
		})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.ServeHTTP(w, r)

			// OPTIONS request don't need any further processing
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type builder struct {
	ContextBuilder ContextBuilder
	Router         *mux.Router
	Prefix         string
	logger         log.Logger
	clock          clock.Clock
}

func (b builder) Add(pattern string, routes server.Routes) {
	url := fmt.Sprintf("%s%s", b.Prefix, pattern)

	methods := make([]string, 0, len(routes))
	for method := range routes {
		methods = append(methods, string(method))
	}

	level.Debug(b.logger).Log("msg", "adding route", "pattern", url)
	b.Router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		iw := &interceptingWriter{code: http.StatusOK, ResponseWriter: w}
		w = iw

		metrics := r.Context().Value(instrumentation.Metrics)
		if registry, ok := metrics.(Instruments); ok {
			counter := registry.Gauge("api")
			counter.Inc()
			defer counter.Dec()

			summary := registry.SummaryVec("api", "method", "path", "code")
			defer func(begin time.Time) {
				summary.WithLabelValues(
					r.Method,
					r.URL.Path,
					strconv.Itoa(iw.code),
				).Observe(b.clock.Since(begin))
			}(b.clock.Now())
		}

		// OPTIONS request don't need any further processing
		if r.Method == "OPTIONS" {
			w.Header().Set("Allow", strings.Join(methods, ", "))
			return
		}

		method := server.Method(strings.ToUpper(r.Method))
		route, ok := routes[method]
		if !ok {
			responses.NotFound(errors.Errorf("no method found for %q", r.URL.String())).Render(b.logger, w)
			return
		}

		ctx, err := b.ContextBuilder.Build(w, r)
		if err != nil {
			responses.InternalError(errors.Wrap(err, "context builder")).Render(b.logger, w)
			return
		}
		resp := route.Func(ctx)
		if resp == nil {
			return
		}
		if err := resp.Render(b.logger, w); err != nil {
			responses.InternalError(err).Render(b.logger, w)
		}
	})
}

type interceptingWriter struct {
	code int
	http.ResponseWriter
}

func (iw *interceptingWriter) WriteHeader(code int) {
	iw.code = code
	iw.ResponseWriter.WriteHeader(code)
}

func (iw *interceptingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := iw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.Errorf("invalid response writer, expected http.Hijacker")
}
