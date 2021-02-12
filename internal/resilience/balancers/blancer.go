package balancers

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/balancers/balancer"
	"github.com/bicycolet/bicycolet/internal/resilience/balancers/random"
	"github.com/bicycolet/bicycolet/internal/resilience/balancers/roundrobin"
	"github.com/pkg/errors"
)

// BalancerType defines the balancer we want to use for resilience.
type BalancerType int

const (
	// RoundRobin balancer that describes what balancer to use.
	RoundRobin BalancerType = iota

	// Random balancer that describes what balancer to use.
	Random
)

// New creates a new encoding balancer based on the type.
func New(t BalancerType, group balancer.Group, options ...Option) (balancer.Balancer, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	switch t {
	case RoundRobin:
		return roundrobin.New(group), nil
	case Random:
		return random.New(group, opts.seed), nil
	default:
		return nil, errors.Errorf("invalid seed type %q", t)
	}
}

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	seed int64
}

// WithSeed sets the seed for token buckets that need a timer.
func WithSeed(seed int64) Option {
	return func(options *options) {
		options.seed = seed
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		seed: time.Now().UnixNano(),
	}
}
