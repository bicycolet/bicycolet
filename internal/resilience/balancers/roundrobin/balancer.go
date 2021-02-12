package roundrobin

import (
	"sync/atomic"

	"github.com/bicycolet/bicycolet/internal/resilience/balancers/balancer"
)

// Balancer creates a round robin load balancer.
type Balancer struct {
	group   balancer.Group
	current uint64
}

// New creates a new round robin balancer.
func New(group balancer.Group) *Balancer {
	return &Balancer{
		group: group,
	}
}

// Index returns a index of the group to use for balancing.
func (b *Balancer) Index() (uint64, error) {
	size := b.group.Size()
	if size <= 0 {
		return 0, balancer.NewErrEmptyGroup()
	}

	old := atomic.AddUint64(&b.current, 1) - 1
	return old % uint64(size), nil
}
