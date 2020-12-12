package component

import (
	"github.com/go-kit/kit/log"
)

// With returns a logger that includes a Key/ComponentValue pair.
func With(logger log.Logger, component string) log.Logger {
	return log.WithPrefix(logger, Key(), ComponentValue(component))
}

// NewFilter wraps next and implements component filtering.
func NewFilter(next log.Logger) log.Logger {
	return &logger{
		next: next,
	}
}

type logger struct {
	next log.Logger
}

func (l *logger) Log(keyvals ...interface{}) error {
	var (
		filtered []interface{}
		found    bool
	)
	for i := 0; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i+1].(*componentValue); ok {
			if found {
				continue
			}
			found = true
		}
		filtered = append(filtered, keyvals[i], keyvals[i+1])
	}
	return l.next.Log(filtered...)
}

// Value is the interface that each of the canonical level values implement.
// It contains unexported methods that prevent types from other packages from
// implementing it and guaranteeing that NewFilter can distinguish the levels
// defined in this package from all other values.
type Value interface {
	String() string
}

// Key returns the unique key added to log events by the loggers in this
// package.
func Key() interface{} { return key }

// ComponentValue returns the unique value added to log events by Component.
func ComponentValue(name string) Value { return &componentValue{name: name} }

// Name returns the underlying component value name.
func Name(component interface{}) string {
	if c, ok := component.(*componentValue); ok {
		return c.String()
	}
	return ""
}

var (
	// key is of type interface{} so that it allocates once during package
	// initialization and avoids allocating every time the value is added to a
	// []interface{} later.
	key interface{} = "component"
)

type componentValue struct {
	name string
}

func (v *componentValue) String() string { return v.name }
