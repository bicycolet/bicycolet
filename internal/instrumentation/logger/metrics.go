package logger

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bicycolet/bicycolet/internal/instrumentation"
)

// Gauge is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
type Gauge struct {
	mutex   sync.RWMutex
	counter int
}

// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
// values.
func (g *Gauge) Inc() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.counter++
}

// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
// values.
func (g *Gauge) Dec() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.counter--
}

// Current returns the underlying current value.
func (g *Gauge) Current() string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	cur := fmt.Sprintf("%d", g.counter)

	g.counter = 0

	return cur
}

// Empty returns if the summary is empty.
func (g *Gauge) Empty() bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.counter == 0
}

// Summary is a Metric that represents a single numerical value that can
//arbitrarily go up and down.
type Summary struct {
	mutex sync.RWMutex
	total int
	mean  float64
	last  int64
}

// Observe adds a single observation to the summary.
func (s *Summary) Observe(d time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.last = int64(d)
	s.total++

	// Rolling average; loss of precision is to be expected. Use a better
	// instrumentation if you want better results.
	s.mean += (float64(d) - s.mean) / float64(s.total)
}

// Current returns the underlying current value.
func (s *Summary) Current() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	cur := fmt.Sprintf("%dt;%v", s.total, time.Duration(s.mean).String())

	s.total = 0
	s.mean = 0
	s.last = 0

	return cur
}

// Empty  returns if the summary is empty.
func (s *Summary) Empty() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.total == 0
}

// SummaryVec is a Collector that bundles a set of Summaries that all share the
// same Desc, but have different values for their variable labels.
type SummaryVec struct {
	mutex     sync.RWMutex
	summaries map[string]*Summary
}

// WithLabelValues returns the Summary for the given slice of label
// values (same order as the VariableLabels in Desc). If that combination of
// label values is accessed for the first time, a new Summary is created.
func (s *SummaryVec) WithLabelValues(labels ...string) instrumentation.Summary {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sort.Strings(labels)
	key := strings.Join(labels, ",")

	if summary, ok := s.summaries[key]; ok {
		return summary
	}

	summary := &Summary{}
	s.summaries[key] = summary
	return summary
}

// Summaries returns the underlying summaries.
func (s *SummaryVec) Summaries() map[string]*Summary {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.summaries
}
