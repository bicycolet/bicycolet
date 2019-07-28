package client

import (
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
)

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	// User agent string
	userAgent string

	// Custom HTTP Client (used as base for the connection)
	httpClient *http.Client

	// Custom proxy
	proxy func(*http.Request) (*url.URL, error)

	// Custom logger
	logger log.Logger
}

// WithUserAgent sets the userAgent on the option
func WithUserAgent(userAgent string) Option {
	return func(options *options) {
		options.userAgent = userAgent
	}
}

// WithHTTPClient sets the httpClient on the option
func WithHTTPClient(httpClient *http.Client) Option {
	return func(options *options) {
		options.httpClient = httpClient
	}
}

// WithProxy sets the proxy on the option
func WithProxy(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(options *options) {
		options.proxy = proxy
	}
}

// WithLogger sets the logger on the option
func WithLogger(logger log.Logger) Option {
	return func(options *options) {
		options.logger = logger
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		httpClient: http.DefaultClient,
		logger:     log.NewNopLogger(),
	}
}
