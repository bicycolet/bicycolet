package breaker

// Breaker creates a new circuit breaker.
type Breaker interface {

	// Run a function against a given breaker.
	Run(func() error) error
}
