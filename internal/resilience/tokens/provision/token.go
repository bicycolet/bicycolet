package provision

import (
	"math"
	"time"

	"github.com/bicycolet/bicycolet/internal/resilience/tickers/ticker"
	"github.com/bicycolet/bicycolet/internal/resilience/tokens/bucket"
	"github.com/bicycolet/bicycolet/internal/resilience/tokens/token"
)

// Token defines a auto provisioning token bucket.
type Token struct {
	bucket *bucket.Token
	ticker ticker.Ticker
}

// New auto provisions a bucket at a given frequency rate
func New(capacity int64, freq time.Duration, ticker token.Ticker) *Token {
	p := &Token{
		bucket: bucket.New(capacity),
	}

	if freq < 0 {
		return p
	} else if evenFreq := time.Duration(1e9 / capacity); freq < evenFreq {
		freq = evenFreq
	}

	inc := int64(math.Floor(.5 + (float64(capacity) * freq.Seconds())))
	p.ticker = ticker.New(freq, func() {
		p.Put(inc)
	})

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
	p.ticker.Stop()
}
