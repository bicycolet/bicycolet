package responses

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Reply represents an API response
type Reply interface {
	Render(log.Logger, http.ResponseWriter) error
}

// ResponseType represents a valid response type
type ResponseType string

// Response types
const (
	Sync  ResponseType = "sync"
	Async ResponseType = "async"
	Error ResponseType = "error"
)

// Response represents a operation
type Response struct {
	Type ResponseType `json:"type" yaml:"type"`

	// Valid only for Sync responses
	Status     string `json:"status" yaml:"status"`
	StatusCode int    `json:"status-code" yaml:"status-code"`

	// Valid only for Async responses
	Operation string `json:"operation" yaml:"operation"`

	// Valid only for Error responses
	Code  int    `json:"error-code" yaml:"error-code"`
	Error string `json:"error" yaml:"error"`

	// Valid for Sync and Error responses
	Metadata json.RawMessage `json:"metadata" yaml:"metadata"`
}

// ResponseRaw represents a operation
type ResponseRaw struct {
	Type ResponseType `json:"type" yaml:"type"`

	// Valid only for Sync responses
	Status     string `json:"status" yaml:"status"`
	StatusCode int    `json:"status-code" yaml:"status-code"`

	// Valid only for Async responses
	Operation string `json:"operation" yaml:"operation"`

	// Valid only for Error responses
	Code  int    `json:"error-code" yaml:"error-code"`
	Error string `json:"error" yaml:"error"`

	// Valid for Sync and Error responses
	Metadata interface{} `json:"metadata" yaml:"metadata"`
}
