package stub

// Ticker represents a stub ticker.
type Ticker struct {
	fn func()
}

// New creates a new ticker based on the interval.
func New(fn func()) *Ticker {
	return &Ticker{
		fn: fn,
	}
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
func (t *Ticker) Stop() {}

// Advance advances the time to call the next tick.
func (t *Ticker) Advance() {
	t.fn()
}
