package bucket

import "sync/atomic"

// Token represents a token bucket implementation.
type Token struct {
	tokens, capacity int64
}

// New creates a new Token with a max capacity
func New(capacity int64) *Token {
	return &Token{
		tokens:   capacity,
		capacity: capacity,
	}
}

// Take attempt to take a number of tokens from a bucket or return the amount
// taken.
func (b *Token) Take(n int64) (taken int64) {
TAKE:
	if tokens := atomic.LoadInt64(&b.tokens); tokens == 0 {
		return 0
	} else if n <= tokens {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, tokens-n) {
			goto TAKE
		}
		return n
	} else {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, 0) {
			goto TAKE
		}
		return tokens
	}
}

// Put attempts to put a number of tokens into a bucket or return the amount
// pushed.
func (b *Token) Put(n int64) (added int64) {
PUT:
	if tokens := atomic.LoadInt64(&b.tokens); tokens == b.capacity {
		return 0
	} else if left := b.capacity - tokens; n <= left {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, tokens+n) {
			goto PUT
		}
		return n
	} else {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, b.capacity) {
			goto PUT
		}
		return left
	}
}

// Close closes a bucket.
func (b *Token) Close() {}
