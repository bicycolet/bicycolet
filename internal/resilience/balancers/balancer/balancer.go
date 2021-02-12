package balancer

// Balancer yields endpoints according to the underlying strategy.
type Balancer interface {
	// Index returns a index of the group to use for balancing.
	Index() (uint64, error)
}

// Group defines a group of blancing items.
type Group interface {
	Size() int
}
