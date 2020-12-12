package noop

import (
	"time"

	"github.com/bicycolet/bicycolet/internal/instrumentation"
)

// Gauge is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
type Gauge struct {
}

// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
// values.
func (g *Gauge) Inc() {}

// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
// values.
func (g *Gauge) Dec() {}

// Summary is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
type Summary struct {
}

// Observe adds a single observation to the summary.
func (s *Summary) Observe(d time.Duration) {}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
type SummaryVec struct {
}

// WithLabelValues returns the Summary for the given slice of label
// values (same order as the VariableLabels in Desc). If that combination of
// label values is accessed for the first time, a new Summary is created.
func (s *SummaryVec) WithLabelValues(labels ...string) instrumentation.Summary {
	return &Summary{}
}
