package client

import (
	"time"

	"github.com/bicycolet/bicycolet/pkg/client"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

// Client describes a very simple client for connecting to a daemon rest API
type Client struct {
	client *client.Client
	logger log.Logger
}

// New creates a Client using the address and certificates.
func New(address string, options ...Option) (*Client, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	client, err := client.New(
		address,
		client.WithLogger(opts.logger),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Client{
		client: client,
		logger: opts.logger,
	}, nil
}

// Get will query the server using the client provided.
func (c *Client) Get(
	path string,
	fn func(*client.Response, Metadata) error,
) error {
	began := time.Now()
	response, responseETag, err := c.client.Query("GET", path, nil, "")
	if err != nil {
		return errors.Wrap(err, "error requesting")
	} else if response.StatusCode != 200 {
		return errors.Errorf("invalid status code %d", response.StatusCode)
	}
	return fn(response, Metadata{
		ETag:     responseETag,
		Duration: time.Since(began),
	})
}

// Metadata holds the metadata for each result.
type Metadata struct {
	ETag     string
	Duration time.Duration
}
