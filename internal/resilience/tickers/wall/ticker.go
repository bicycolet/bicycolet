package wall

import (
	"time"
)

// Ticker represents a wall ticker.
type Ticker struct {
	stop     chan chan struct{}
	ticker   *time.Ticker
	interval time.Duration
}

// New creates a new ticker based on the interval.
func New(interval time.Duration, fn func()) *Ticker {
	ticker := time.NewTicker(interval)
	t := &Ticker{
		stop:     make(chan chan struct{}),
		interval: interval,
		ticker:   ticker,
	}

	go func() {
		for {
			select {
			case c := <-t.stop:
				close(c)
				return
			case <-ticker.C:
				fn()
			}
		}
	}()

	return t
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
func (t *Ticker) Stop() {
	t.ticker.Stop()

	c := make(chan struct{})
	t.stop <- c
	<-c
}
