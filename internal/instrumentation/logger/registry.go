package logger

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"

	"github.com/bicycolet/bicycolet/internal/instrumentation"
	"github.com/go-kit/kit/log"
	"github.com/spoke-d/task"
	"github.com/spoke-d/task/wait"
)

const (
	// Interval of how often the registry should print out the metrics
	Interval = time.Second * 5
)

// Registry defines a logger instrumentation.
type Registry struct {
	mutex       sync.RWMutex
	gauges      map[string]*Gauge
	summaries   map[string]*Summary
	summaryVecs map[string]*SummaryVec
	logger      log.Logger
}

// New creates a new Metrics registry.
func New(logger log.Logger) *Registry {
	return &Registry{
		gauges:      make(map[string]*Gauge),
		summaries:   make(map[string]*Summary),
		summaryVecs: make(map[string]*SummaryVec),
		logger:      logger,
	}
}

// Run returns a task function that performs the outputting of the metrics
func (r *Registry) Run() (task.Func, task.Schedule) {
	schedulerWrapper := func(ctx context.Context) error {
		return wait.Wait(r.run, time.Second*30, wait.WithContext(ctx))
	}

	schedule := task.Every(Interval)
	return schedulerWrapper, schedule
}

// Gauge is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
func (r *Registry) Gauge(name string) instrumentation.Gauge {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if g, ok := r.gauges[name]; ok {
		return g
	}

	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Summary captures individual observations from an event or sample stream and
// summarizes them in a manner similar to traditional summary statistics.
func (r *Registry) Summary(name string) instrumentation.Summary {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if s, ok := r.summaries[name]; ok {
		return s
	}

	s := &Summary{}
	r.summaries[name] = s
	return s
}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
func (r *Registry) SummaryVec(name string, labels ...string) instrumentation.SummaryVec {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if s, ok := r.summaryVecs[name]; ok {
		return s
	}

	s := &SummaryVec{
		summaries: make(map[string]*Summary),
	}
	r.summaryVecs[name] = s
	return s
}

func (r *Registry) run(ctx context.Context) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for name, gauge := range r.gauges {
		if gauge.Empty() {
			continue
		}
		level.Debug(r.logger).Log("metric", "gauge", "name", name, "current", gauge.Current())
	}

	for name, summary := range r.summaries {
		if summary.Empty() {
			continue
		}
		level.Debug(r.logger).Log("metric", "summary", "name", name, "current", summary.Current())
	}

	for name, vec := range r.summaryVecs {
		summaries := vec.Summaries()
		for key, summary := range summaries {
			if summary.Empty() {
				continue
			}
			level.Debug(r.logger).Log("metric", "summaryvec", "name", name, "key", key, "current", summary.Current())
		}
	}
}
