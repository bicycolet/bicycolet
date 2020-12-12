package middleware

import (
	"context"
	"net/http"
)

// ContextBuilder will aim to build the context.
type ContextBuilder interface {
	// Build the context passing in the http types.
	Build(http.ResponseWriter, *http.Request) (context.Context, error)
}
