// Package gob provides a codec for encoding and decoding gob data.
package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/bicycolet/bicycolet/internal/codecs/encoding"
)

// Codec creates a codec for encoding and decoding gob data.
type Codec struct{}

// New returns a codec that implements encoding and decoding of gob data.
func New() *Codec {
	return &Codec{}
}

// Marshaler is the interface implemented by objects that can produce a byte
// representation of another object.
func (c *Codec) Marshaler() encoding.Marshaler {
	return func(value interface{}) ([]byte, error) {
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(value)
		return buf.Bytes(), err
	}
}

// Unmarshaler is the interface implemented by objects that can unmarshal a byte
// representation of an object into an object instance.
func (c *Codec) Unmarshaler() encoding.Unmarshaler {
	return func(data []byte, target interface{}) error {
		return gob.NewDecoder(bytes.NewReader(data)).Decode(target)
	}
}
