package demux

import "github.com/go-kit/kit/log"

// Logger defines a new Logger for use with convoy.
type Logger struct {
	debug         bool
	syslog, other log.Logger
}

// NewLogger creates a new Logger for use with convoy.
func NewLogger(debug bool, syslog, other log.Logger) *Logger {
	return &Logger{
		debug:  debug,
		syslog: syslog,
		other:  other,
	}
}

// Log out key value pairs.
func (c *Logger) Log(keyvals ...interface{}) error {
	// If in debug mode, don't send to syslog
	if !c.debug {
		if err := c.syslog.Log(keyvals...); err != nil {
			panic(err)
		}
	}

	if err := c.other.Log(keyvals...); err != nil {
		panic(err)
	}
	return nil
}
