package throttle

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tsenart/tb"
)

// NewFilter wraps next and implements component filtering.
func NewFilter(next log.Logger, key string, amount int64, freq time.Duration) log.Logger {
	return &logger{
		next:   next,
		key:    key,
		bucket: tb.NewBucket(amount, freq),
	}
}

type logger struct {
	next   log.Logger
	key    interface{}
	bucket *tb.Bucket
}

func (l *logger) Log(keyvals ...interface{}) error {
	var filtered []interface{}
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == l.key && l.bucket.Take(1) == 0 {
			return nil
		}
		filtered = append(filtered, keyvals[i], keyvals[i+1])
	}
	return l.next.Log(filtered...)
}
