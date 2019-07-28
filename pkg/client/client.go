package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Client represents interactions with a API server
type Client struct {
	http          *http.Client
	httpHost      string
	httpProtocol  string
	httpUserAgent string

	logger log.Logger
}

// New lets you connect to a remote daemon over HTTP.
func New(url string, options ...Option) (*Client, error) {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	// Setup the HTTP client
	httpClient, err := tlsHTTPClient(opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Initialize the client struct
	return &Client{
		httpHost:      url,
		httpProtocol:  "http",
		httpUserAgent: opts.userAgent,
		http:          httpClient,
		logger:        opts.logger,
	}, nil
}

// HTTPClient returns the http client used for the connection.
// This can be used to set custom http options.
func (c *Client) HTTPClient() *http.Client {
	return c.http
}

// Query allows directly querying the API
func (c *Client) Query(
	method, path string,
	data interface{},
	eTag string,
) (*Response, string, error) {
	// Generate the URL
	url := fmt.Sprintf("%s%s", c.httpHost, path)
	return c.rawQuery(method, url, data, eTag)
}

// Websocket allows directly connection to API websockets
func (c *Client) Websocket(path string) (*websocket.Conn, error) {
	return c.websocket(path)
}

// Do performs a Request.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

func (c *Client) rawQuery(method, url string, data interface{}, eTag string) (*Response, string, error) {
	var req *http.Request
	var err error

	// Log the request
	level.Debug(c.logger).Log("msg", "sending request to API",
		"method", method,
		"url", url,
		"etag", eTag,
	)

	// Get a new HTTP request setup
	if data != nil {
		switch data.(type) {
		case io.Reader:
			// Some data to be sent along with the request
			req, err = http.NewRequest(method, url, data.(io.Reader))
			if err != nil {
				return nil, "", errors.WithStack(err)
			}

			// Set the encoding accordingly
			req.Header.Set("Content-Type", "application/octet-stream")
		default:
			// Encode the provided data
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(data); err != nil {
				return nil, "", errors.WithStack(err)
			}

			// Some data to be sent along with the request
			// Use a reader since the request body needs to be seekable
			req, err = http.NewRequest(method, url, bytes.NewReader(buf.Bytes()))
			if err != nil {
				return nil, "", errors.WithStack(err)
			}

			// Set the encoding accordingly
			req.Header.Set("Content-Type", "application/json")

			// Log the data
			level.Debug(c.logger).Log("data", data)
		}
	} else {
		// No data to be sent along with the request
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, "", errors.WithStack(err)
		}
	}

	// Set the user agent
	if c.httpUserAgent != "" {
		req.Header.Set("User-Agent", c.httpUserAgent)
	}

	// Set the eTag
	if eTag != "" {
		req.Header.Set("If-Match", eTag)
	}

	// Send the request
	resp, err := c.do(req)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

func (c *Client) websocket(path string) (*websocket.Conn, error) {
	host := c.httpHost
	// Generate the URL
	var url string
	if strings.HasPrefix(host, "https://") {
		url = fmt.Sprintf("wss://%s/1.0%s", strings.TrimPrefix(host, "https://"), path)
	} else {
		url = fmt.Sprintf("ws://%s/1.0%s", strings.TrimPrefix(host, "http://"), path)
	}

	return c.rawWebsocket(url)
}

func (c *Client) rawWebsocket(url string) (*websocket.Conn, error) {
	// Grab the http transport handler
	httpTransport, ok := c.http.Transport.(*http.Transport)
	if !ok {
		return nil, errors.New("invalid transport")
	}

	// Setup a new websocket dialer based on it
	dialer := websocket.Dialer{
		NetDial:         httpTransport.Dial,
		TLSClientConfig: httpTransport.TLSClientConfig,
		Proxy:           httpTransport.Proxy,
	}

	// Set the user agent
	headers := http.Header{}
	if c.httpUserAgent != "" {
		headers.Set("User-Agent", c.httpUserAgent)
	}

	// Establish the connection
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Log the data
	level.Debug(c.logger).Log("msg", "Connected to the websocket")

	return conn, err
}

func tlsHTTPClient(opts *options) (*http.Client, error) {
	// Define the http transport
	transport := &http.Transport{}

	// Allow overriding the proxy
	proxy := opts.proxy
	if proxy != nil {
		transport.Proxy = proxy
	}

	// Define the http client
	client := opts.httpClient
	if client == nil {
		client = &http.Client{}
	}
	client.Transport = transport

	// Setup redirect policy
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Replicate the headers
		req.Header = via[len(via)-1].Header
		return nil
	}

	return client, nil
}

func parseResponse(resp *http.Response) (*Response, string, error) {
	// Get the ETag
	etag := resp.Header.Get("ETag")

	var response Response
	// Empty response body
	if resp.ContentLength == 0 {
		response.Type = SyncResponse
		response.Status = resp.Status
		response.StatusCode = resp.StatusCode
		return &response, etag, nil
	}

	// Decode the response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		// Check the return value for a cleaner error
		if resp.StatusCode != http.StatusOK {
			return nil, "", errors.Errorf("failed to fetch %s: %s", resp.Request.URL.String(), resp.Status)
		}

		return nil, "", errors.WithStack(err)
	}

	// Handle errors
	if response.Type == ErrorResponse {
		return nil, "", errors.New(response.Error)
	}

	return &response, etag, nil
}
