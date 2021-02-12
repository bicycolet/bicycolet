package breaker

// errBreakerOpen defines an err breaker.
type errBreakerOpen struct{}

// NewErrBreakerOpen creates a new break is open error
func NewErrBreakerOpen() error {
	return errBreakerOpen{}
}

func (e errBreakerOpen) Error() string {
	return "breaker is open"
}

// ErrBreakerOpen checks if the error was because of the breaker being in the
// open state.
func ErrBreakerOpen(err error) bool {
	_, ok := err.(errBreakerOpen)
	return ok
}
