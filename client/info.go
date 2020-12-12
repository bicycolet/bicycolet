package client

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/bicycolet/bicycolet/api/handlers/info"
	"github.com/bicycolet/bicycolet/api/responses"
	"github.com/pkg/errors"
)

// InfoAPI represents a way of interacting with the daemon API, which is
// responsible for getting the information from the daemon.
type InfoAPI struct {
	client *Client
}

// NewInfoAPI creates a API client for requesting server information
func NewInfoAPI(client *Client) *InfoAPI {
	return &InfoAPI{
		client: client,
	}
}

// Get returns the information from the daemon API
func (i *InfoAPI) Get() (InfoResult, string, error) {
	var etag string
	var result InfoResult
	if err := i.client.exec(context.TODO(), "GET", "/1.0", nil, "", func(response *responses.Response, meta Metadata) error {
		var server info.Envelope
		reader := bytes.NewReader(response.Metadata)
		if err := json.NewDecoder(reader).Decode(&server); err != nil {
			return errors.Wrap(err, "error parsing result")
		}

		config := make(map[string]interface{}, len(server.Config))
		for k, v := range server.Config {
			config[k] = v
		}

		etag = meta.ETag
		result = InfoResult{
			Environment: Environment{
				Addresses:              server.Environment.Addresses,
				Certificate:            server.Environment.Certificate,
				CertificateFingerprint: server.Environment.CertificateFingerprint,
				CertificateKey:         server.Environment.CertificateKey,
				Server:                 server.Environment.Server,
				ServerPid:              server.Environment.ServerPid,
				ServerVersion:          server.Environment.ServerVersion,
				ServerClustered:        server.Environment.ServerClustered,
				ServerName:             server.Environment.ServerName,
			},
			Config: config,
		}
		return nil
	}); err != nil {
		return result, "", errors.WithStack(err)
	}
	return result, etag, nil
}

// InfoResult contains the result of querying the daemon information API
type InfoResult struct {
	Environment Environment            `json:"environment" yaml:"environment"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
}

// Environment defines the server environment for the daemon
type Environment struct {
	Addresses              []string `json:"addresses" yaml:"addresses"`
	Certificate            string   `json:"certificate" yaml:"certificate"`
	CertificateFingerprint string   `json:"certificate-fingerprint" yaml:"certificate-fingerprint"`
	CertificateKey         string   `json:"certificate-key,omitempty" yaml:"certificate-key,omitempty"`
	Server                 string   `json:"server" yaml:"server"`
	ServerPid              int      `json:"server-pid" yaml:"server-pid"`
	ServerVersion          string   `json:"server-version" yaml:"server-version"`
	ServerClustered        bool     `json:"server-clustered" yaml:"server-clustered"`
	ServerName             string   `json:"server-name" yaml:"server-name"`
}

// Map returns a map[string]interface{} of the environment setup
func (e Environment) Map() map[string]interface{} {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(e)

	var result map[string]interface{}
	json.NewDecoder(&buf).Decode(&result)

	return result
}
