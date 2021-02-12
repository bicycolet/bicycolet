package random

import (
	"math/rand"

	"github.com/bicycolet/bicycolet/internal/resilience/balancers/balancer"
)

// Balancer creates a round robin load balancer.
type Balancer struct {
	group balancer.Group
	rand  *rand.Rand
}

// New creates a new round robin balancer.
func New(group balancer.Group, seed int64) *Balancer {
	return &Balancer{
		group: group,
		rand:  rand.New(rand.NewSource(seed)),
	}
}

// Index returns a index of the group to use for balancing.
func (b *Balancer) Index() (uint64, error) {
	size := b.group.Size()
	if size <= 0 {
		return 0, balancer.NewErrEmptyGroup()
	}

	return uint64(b.rand.Intn(size)), nil
}
