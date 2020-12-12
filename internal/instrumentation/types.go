package instrumentation

import "time"

// Gauge is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
//
// A Gauge is typically used for measured values like temperatures or current
// memory usage, but also "counts" that can go up and down, like the number of
// running goroutines.
type Gauge interface {
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc()
	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
	// values.
	Dec()
}

// Summary captures individual observations from an event or sample stream and
// summarizes them in a manner similar to traditional summary statistics:
//  1. sum of observations
//  2. observation count
//  3. rank estimations.
//
// A typical use-case is the observation of request latencies. By default, a
// Summary provides the median, the 90th and the 99th percentile of the latency
// as rank estimations. However, the default behavior will change in the
// upcoming v1.0.0 of the library. There will be no rank estimations at all by
// default. For a sane transition, it is recommended to set the desired rank
// estimations explicitly.
type Summary interface {
	// Observe adds a single observation to the summary.
	Observe(time.Duration)
}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
type SummaryVec interface {
	// WithLabelValues returns the Summary for the given slice of label
	// values (same order as the VariableLabels in Desc). If that combination of
	// label values is accessed for the first time, a new Summary is created.
	WithLabelValues(...string) Summary
}
