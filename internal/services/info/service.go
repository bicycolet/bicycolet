package info

import (
	"context"

	"github.com/bicycolet/bicycolet/internal/services"
)

const ServiceKey services.ServiceKey = "info"

const (
	// ServerName is the name of default server in non-clustered mode.
	ServerName = "bicycolet"
)

// Envelope represents the structure for the server
type Envelope struct {
	Environment Environment
}

// Environment defines the server environment for the daemon
type Environment struct {
	Addresses     []string
	Server        string
	ServerPid     int
	ServerVersion string
	ServerName    string
}

// Service represents a service for getting and requesting data.
type Service struct {
}

// New creates a new Info service
func New(options ...Option) *Service {
	opts := newOptions()
	for _, option := range options {
		option(opts)
	}

	return &Service{}
}

// Get defines a service for calling "GET" method and returns a response.
func (s *Service) Get(ctx context.Context) (Envelope, error) {
	return Envelope{
		Environment: Environment{
			ServerName: ServerName,
		},
	}, nil
}
