package provision

import (
	"math"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tokens/bucket"
)

// Token defines a auto provisioning token bucket.
type Token struct {
	bucket *bucket.Token
	freq   time.Duration
	inc    int64
	stop   chan chan struct{}
}

// New auto provisions a bucket at a given frequency rate
func New(capacity int64, freq time.Duration) *Token {
	p := &Token{
		bucket: bucket.New(capacity),
		freq:   freq,
		stop:   make(chan chan struct{}),
	}

	if freq < 0 {
		return p
	} else if evenFreq := time.Duration(1e9 / capacity); freq < evenFreq {
		freq = evenFreq
	}

	p.freq = freq
	p.inc = int64(math.Floor(.5 + (float64(capacity) * freq.Seconds())))

	return p
}

// Take attempts to take n tokens out of the bucket.
func (p *Token) Take(n int64) int64 {
	return p.bucket.Take(n)
}

// Put attempts to add n tokens to the bucket.
func (p *Token) Put(n int64) int64 {
	return p.bucket.Put(n)
}

// Close stops the filling of a given bucket if it was started.
func (p *Token) Close() {
	c := make(chan struct{})
	p.stop <- c
	<-c
}

// Run the token to ensure that the bucket is updated over a given frequency.
func (p *Token) Run() {
	ticker := time.NewTicker(p.freq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.Put(p.inc)
		case q := <-p.stop:
			close(q)
			return
		}
	}
}
