package encoding

// Marshaler is the interface implemented by objects that can produce a byte
// representation of another object.
type Marshaler func(interface{}) ([]byte, error)

// Unmarshaler is the interface implemented by objects that can unmarshal a byte
// representation of an object into an object instance.
type Unmarshaler func([]byte, interface{}) error

// Codec is implemented by objects that can produce marshalers and unmarshalers
// for a given object type.
//
// Marshaler returns a Marshaler implementation that can marshal instances of a
// particular type into a byte slice.
//
// Unmarshaler returns a Unmarshaler implementation that can unmarshal instances
// of a particular type from a byte slice.
type Codec interface {
	Marshaler() Marshaler
	Unmarshaler() Unmarshaler
}
