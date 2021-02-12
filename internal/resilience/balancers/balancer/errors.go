package balancer

// errEmptyGroup defines an err balancer.
type errEmptyGroup struct{}

// NewErrEmptyGroup creates a new error describing empty groups.
func NewErrEmptyGroup() error {
	return errEmptyGroup{}
}

func (e errEmptyGroup) Error() string {
	return "empty group"
}

// ErrBreakerOpen checks if the error was because of the balancer being in the
// open state.
func ErrBreakerOpen(err error) bool {
	_, ok := err.(errEmptyGroup)
	return ok
}
