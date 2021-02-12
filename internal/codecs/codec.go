package codecs

import (
	"github.com/bicycolet/bicycolet/internal/codecs/encoding"
	"github.com/bicycolet/bicycolet/internal/codecs/gob"
	"github.com/bicycolet/bicycolet/internal/codecs/json"
	"github.com/pkg/errors"
)

// CodecType defines the codec we want to use for transporting.
type CodecType int

const (
	// JSON is a codec that marshals and unmarshals data in the JSON format.
	JSON CodecType = iota
	// GOB is a codec that marshals and unmarshals data in the gob format.
	GOB
)

// New creates a new encoding codec based on the type.
func New(t CodecType) (encoding.Codec, error) {
	switch t {
	case JSON:
		return json.New(), nil
	case GOB:
		return gob.New(), nil
	default:
		return nil, errors.Errorf("invalid codec type %q", t)
	}
}
