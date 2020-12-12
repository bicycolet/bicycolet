package responses

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
)

// SyncResponse represents a sync response.
type SyncResponse struct {
	success  bool
	eTag     interface{}
	metadata interface{}
	location string
	code     int
	headers  map[string]string
}

// Render will consume a http.ResponseWriter and return an error in a vistor
// pattern scenario.
func (r *SyncResponse) Render(log log.Logger, w http.ResponseWriter) error {
	// Set an appropriate ETag header
	if r.eTag != nil {
		if eTag, err := Hash(r.eTag); err == nil {
			w.Header().Set("ETag", eTag)
		}
	}

	status := http.StatusOK
	if !r.success {
		status = http.StatusBadRequest
	}

	if r.headers != nil {
		for h, v := range r.headers {
			w.Header().Set(h, v)
		}
	}

	if r.location != "" {
		w.Header().Set("Location", r.location)
		if r.code == 0 {
			w.WriteHeader(201)
		} else {
			w.WriteHeader(r.code)
		}
	}

	return json.NewEncoder(w).Encode(ResponseRaw{
		Type:       Sync,
		Status:     http.StatusText(status),
		StatusCode: status,
		Metadata:   r.metadata,
	})
}

// SyncResponsePermanentRedirect defines a successful response that will always perform a
// permanent redirect.
func SyncResponsePermanentRedirect(address string) *SyncResponse {
	return &SyncResponse{
		success:  true,
		location: address,
		code:     http.StatusPermanentRedirect,
	}
}

// SyncResponseTemporaryRedirect defines a successful response that will always perform a
// permanent redirect.
func SyncResponseTemporaryRedirect(address string) *SyncResponse {
	return &SyncResponse{
		success:  true,
		location: address,
		code:     http.StatusTemporaryRedirect,
	}
}

// SyncResponseETag defines a response that can add ETag as additional
// information
func SyncResponseETag(success bool, metadata interface{}, eTag interface{}) *SyncResponse {
	return &SyncResponse{
		success:  success,
		metadata: metadata,
		eTag:     eTag,
	}
}

// EmptySyncResponse defines an empty successful response
func EmptySyncResponse() *SyncResponse {
	return &SyncResponse{
		success:  true,
		metadata: make(map[string]interface{}),
	}
}
