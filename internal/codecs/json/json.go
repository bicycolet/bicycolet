// Package json provides a codec for encoding and decoding JSON data.
package json

import (
	"encoding/json"

	"github.com/bicycolet/bicycolet/internal/codecs/encoding"
)

// Option to be passed to Connect to customize the resulting instance.
type Option func(*options)

type options struct {
	pretty bool
}

// WithPrettyOutput sets pretty output for the codec.
func WithPrettyOutput(pretty bool) Option {
	return func(options *options) {
		options.pretty = pretty
	}
}

// Create a options instance with default values.
func newOptions() *options {
	return &options{
		pretty: false,
	}
}

// Codec creates a codec for encoding and decoding JSON data.
type Codec struct {
	pretty bool
}

// New returns a codec that implements encoding and decoding of JSON data.
func New(options ...Option) *Codec {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	return &Codec{
		pretty: opts.pretty,
	}
}

// Marshaler is the interface implemented by objects that can produce a byte
// representation of another object.
func (c *Codec) Marshaler() encoding.Marshaler {
	if c.pretty {
		return func(value interface{}) ([]byte, error) {
			return json.MarshalIndent(value, "", "  ")
		}
	}
	return json.Marshal
}

// Unmarshaler is the interface implemented by objects that can unmarshal a byte
// representation of an object into an object instance.
func (c *Codec) Unmarshaler() encoding.Unmarshaler {
	return json.Unmarshal
}
