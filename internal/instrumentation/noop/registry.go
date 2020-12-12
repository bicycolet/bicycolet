package noop

import (
	"github.com/bicycolet/bicycolet/internal/instrumentation"
)

// Registry defines a logger instrumentation.
type Registry struct {
	gauge      *Gauge
	summary    *Summary
	summaryVec *SummaryVec
}

// New creates a new Metrics registry.
func New() *Registry {
	return &Registry{
		gauge:      &Gauge{},
		summary:    &Summary{},
		summaryVec: &SummaryVec{},
	}
}

// Gauge is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
func (r *Registry) Gauge(name string) instrumentation.Gauge {
	return r.gauge
}

// Summary captures individual observations from an event or sample stream and
// summarizes them in a manner similar to traditional summary statistics.
func (r *Registry) Summary(name string) instrumentation.Summary {
	return r.summary
}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
func (r *Registry) SummaryVec(name string, labels ...string) instrumentation.SummaryVec {
	return r.summaryVec
}
