package info

import (
	"bytes"
	"encoding/json"

	"github.com/bicycolet/bicycolet/pkg/api/daemon/root"
	"github.com/pkg/errors"
)

// Response represents a response from the server
type Response interface {
	Metadata() json.RawMessage
}

// Metadata represents the data associated with a request/response
type Metadata interface{}

// ResponseFn is an alias for a response callback
type ResponseFn = func(Response, Metadata) error

// Client represents a point of use root client
type Client interface {
	Get(url string, response ResponseFn) error
}

// Info represents a way of interacting with the daemon API, which is
// responsible for getting the information from the daemon.
type Info struct {
	client Client
}

// New creates a new info client for dealing with information requests.
func New(client Client) *Info {
	return &Info{
		client: client,
	}
}

// Get returns the information from the daemon API
func (i Info) Get() (InfoResult, error) {
	var result InfoResult
	if err := i.client.Get("/1.0", func(response Response, meta Metadata) error {
		var server root.Server
		decoder := json.NewDecoder(bytes.NewReader(response.Metadata()))
		if err := decoder.Decode(&server); err != nil {
			return errors.Wrap(err, "error parsing result")
		}

		config := make(map[string]interface{}, len(server.Config))
		for k, v := range server.Config {
			config[k] = v
		}

		result = InfoResult{
			Environment: Environment{
				Addresses:     server.Environment.Addresses,
				Server:        server.Environment.Server,
				ServerPid:     server.Environment.ServerPid,
				ServerVersion: server.Environment.ServerVersion,
				ServerName:    server.Environment.ServerName,
			},
			Config: config,
		}
		return nil
	}); err != nil {
		return result, errors.WithStack(err)
	}
	return result, nil
}

// InfoResult contains the result of querying the daemon information API
type InfoResult struct {
	Environment Environment            `json:"environment" yaml:"environment"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
}

// Environment defines the server environment for the daemon
type Environment struct {
	Addresses     []string `json:"addresses" yaml:"addresses"`
	Server        string   `json:"server" yaml:"server"`
	ServerPid     int      `json:"server_pid" yaml:"server_pid"`
	ServerVersion string   `json:"server_version" yaml:"server_version"`
	ServerName    string   `json:"server_name" yaml:"server_name"`
}
