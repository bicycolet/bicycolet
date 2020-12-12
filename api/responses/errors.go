package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// ErrorResponse encapsulates an api response.
type ErrorResponse struct {
	code   int
	msg    string
	origin error
}

func (r *ErrorResponse) String() string {
	return r.msg
}

// Render will consume a http.ResponseWriter and return an error in a vistor
// pattern scenario.
func (r *ErrorResponse) Render(logger log.Logger, w http.ResponseWriter) error {
	if r.origin != nil && r.code >= 400 {
		level.Debug(logger).Log("code", r.code, "error", fmt.Sprintf("%v", r.origin))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(r.code)

	return json.NewEncoder(w).Encode(map[string]interface{}{
		"type":       Error,
		"error":      r.msg,
		"error-code": r.code,
	})
}

// NotImplemented takes an error and returns a Response of not implemented.
func NotImplemented(err error) *ErrorResponse {
	message := "not implemented"
	if err != nil {
		message = err.Error()
	}
	return &ErrorResponse{
		code: http.StatusNotImplemented,
		msg:  message,
	}
}

// NotFound takes an error and returns a Response of not found.
func NotFound(err error) *ErrorResponse {
	message := "not found"
	if err != nil {
		message = err.Error()
	}
	return &ErrorResponse{
		code: http.StatusNotFound,
		msg:  message,
	}
}

// Forbidden takes an error and returns a Response of forbidden error.
func Forbidden(err error) *ErrorResponse {
	message := "not authorized"
	if err != nil {
		message = err.Error()
	}
	return &ErrorResponse{
		code:   http.StatusForbidden,
		msg:    message,
		origin: err,
	}
}

// Conflict takes an error and returns a Response of conflict error.
func Conflict(err error) *ErrorResponse {
	message := "already exists"
	if err != nil {
		message = err.Error()
	}
	return &ErrorResponse{
		code:   http.StatusConflict,
		msg:    message,
		origin: err,
	}
}

// Unavailable takes an error and returns a Response of unavailable error.
func Unavailable(err error) *ErrorResponse {
	message := "unavailable"
	if err != nil {
		message = err.Error()
	}
	return &ErrorResponse{
		code:   http.StatusServiceUnavailable,
		msg:    message,
		origin: err,
	}
}

// BadRequest takes an error and returns a Response of badrequest error.
func BadRequest(err error) *ErrorResponse {
	return &ErrorResponse{
		code:   http.StatusBadRequest,
		msg:    err.Error(),
		origin: err,
	}
}

// InternalError takes an error and returns a Response of internal server error.
func InternalError(err error) *ErrorResponse {
	return &ErrorResponse{
		code:   http.StatusInternalServerError,
		msg:    err.Error(),
		origin: err,
	}
}

// PreconditionFailed takes an error and returns a Response of precondition
// failed error.
func PreconditionFailed(err error) *ErrorResponse {
	return &ErrorResponse{
		code:   http.StatusPreconditionFailed,
		msg:    err.Error(),
		origin: err,
	}
}
