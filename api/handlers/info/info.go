package info

import (
	"context"

	"github.com/bicycolet/bicycolet/api/middleware/server"
	"github.com/bicycolet/bicycolet/api/responses"
	"github.com/bicycolet/bicycolet/internal/services/info"
	"github.com/pkg/errors"
)

// Router defines how the handler subscribes to the available routes on offer.
type Router = server.Router

// Routes defines a mapping of methods to routes.
type Routes = server.Routes

// Route is a route that is triggered when a method and pattern is matched.
type Route = server.Route

// Handler represents a Info handler.
type Handler struct{}

// New creates a new ready Handler.
func New() *Handler {
	return &Handler{}
}

// Subscribe attempts to subscribe to the offering endpoints to the Router.
func (h *Handler) Subscribe(router Router) {
	router.Add("", server.Single(server.Get, server.Unrestricted(h.GetInfo)))
}

// GetInfo executes if the `/ready` endpoint is triggered. It workers
func (h *Handler) GetInfo(ctx context.Context) server.Response {
	service := ctx.Value(info.ServiceKey)
	if service == nil {
		return responses.InternalError(errors.Errorf("info service not found"))
	}
	infoService, ok := service.(InfoService)
	if !ok {
		return responses.InternalError(errors.Errorf("info service not valid"))
	}

	// Perform request.

	info, err := infoService.Get(ctx)
	if err != nil {
		return responses.InternalError(err)
	}
	env := info.Environment
	envelope := Envelope{
		Environment: Environment{
			ServerName: env.ServerName,
		},
	}

	return responses.SyncResponseETag(true, envelope, envelope.Config)
}

// InfoEnvelope is the data returned from the info service.
type InfoEnvelope = info.Envelope

// InfoService handles all info requests from the server.
type InfoService interface {

	// Get returns the service information including the config.
	Get(context.Context) (InfoEnvelope, error)
}
