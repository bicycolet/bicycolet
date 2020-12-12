package prometheus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bicycolet/bicycolet/internal/instrumentation"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spoke-d/task"
)

// Registerer defines a metric Registerer from prometheus
type Registerer interface {
	prometheus.Registerer
	prometheus.Gatherer
}

const (
	// Interval of how often the registry should print out the metrics
	Interval = time.Second * 60
)

// Registry defines a logger instrumentation.
type Registry struct {
	mutex       sync.RWMutex
	registry    Registerer
	gauges      map[string]prometheus.Gauge
	summaries   map[string]prometheus.Summary
	summaryVecs map[string]*prometheus.SummaryVec
	logger      log.Logger
}

// New creates a new Metrics registry.
func New(registry Registerer, logger log.Logger) *Registry {
	return &Registry{
		registry:    registry,
		gauges:      make(map[string]prometheus.Gauge),
		summaries:   make(map[string]prometheus.Summary),
		summaryVecs: make(map[string]*prometheus.SummaryVec),
		logger:      logger,
	}
}

// Run returns a task function that performs the outputting of the metrics
func (r *Registry) Run() (task.Func, task.Schedule) {
	schedulerWrapper := func(ctx context.Context) error {
		return nil
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

	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: name,
		Name:      fmt.Sprintf("%s_guage", name),
		Help:      fmt.Sprintf("How many %s calls there are", name),
	})
	if err := r.registry.Register(g); err != nil {
		level.Error(r.logger).Log("msg", "failed to register gauge", "name", name)
		return g
	}
	r.gauges[name] = g
	return g
}

// Summary captures individual observations from an event or sample stream and
// summarizes them in a manner similar to traditional summary statistics.
func (r *Registry) Summary(name string) instrumentation.Summary {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if s, ok := r.summaries[name]; ok {
		return summary{Summary: s}
	}

	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: name,
		Name:      fmt.Sprintf("%s_summary", name),
		Help:      fmt.Sprintf("How long the %s calls took", name),
	})
	if err := r.registry.Register(s); err != nil {
		level.Error(r.logger).Log("msg", "failed to register summary", "name", name)
		return summary{Summary: s}
	}
	r.summaries[name] = s
	return summary{Summary: s}
}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
func (r *Registry) SummaryVec(name string, labels ...string) instrumentation.SummaryVec {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if s, ok := r.summaryVecs[name]; ok {
		return summaryVec{SummaryVec: s}
	}

	s := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: name,
		Name:      fmt.Sprintf("%s_summary_vec", name),
		Help:      fmt.Sprintf("How long the %s calls took", name),
	}, labels)
	if err := r.registry.Register(s); err != nil {
		level.Error(r.logger).Log("msg", "failed to register summary vec", "name", name, "err", err)
		return summaryVec{SummaryVec: s}
	}
	r.summaryVecs[name] = s
	return summaryVec{SummaryVec: s}
}

type summary struct {
	prometheus.Summary
}

func (s summary) Observe(d time.Duration) {
	s.Summary.Observe(d.Seconds())
}

type summaryVec struct {
	*prometheus.SummaryVec
}

func (s summaryVec) WithLabelValues(labels ...string) instrumentation.Summary {
	return summaryVecObserver{
		Observer: s.SummaryVec.WithLabelValues(labels...),
	}
}

type summaryVecObserver struct {
	prometheus.Observer
}

func (s summaryVecObserver) Observe(d time.Duration) {
	s.Observer.Observe(d.Seconds())
}
