package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Method defines which method to add to the router.
type Method string

const (
	Get    Method = "GET"
	Put    Method = "PUT"
	Post   Method = "POST"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

// Response is returns back from a route to be rendered.
type Response interface {
	Render(log.Logger, http.ResponseWriter) error
}

// Single constructs a singular routing table for a Method route.
func Single(m Method, r Route) Routes {
	return map[Method]Route{
		m: r,
	}
}

// Router allows the adding of a new routable endpoint.
type Router interface {
	// Add defines a pattern and a callback when the pattern is matched.
	Add(string, Routes)
}

// Routes defines a mapping of methods to routes.
type Routes map[Method]Route

// Func defines a route function
type Func = func(context.Context) Response

// Route is a route that is triggered when a method and pattern is matched.
type Route struct {
	Func                   Func
	AllowedUnauthenticated bool
	AllowedUntrusted       bool
	ValidPath              bool
}

// Trusted defines a trusted route.
func Trusted(fn Func) Route {
	return Route{
		Func:                   fn,
		AllowedUnauthenticated: true,
		AllowedUntrusted:       false,
		ValidPath:              false,
	}
}

// Authenticated defines a trusted route.
func Authenticated(fn Func) Route {
	return Route{
		Func:                   fn,
		AllowedUnauthenticated: false,
		AllowedUntrusted:       true,
		ValidPath:              false,
	}
}

// AuthenticatedOrTrusted defines a trusted route.
func AuthenticatedOrTrusted(fn Func) Route {
	return Route{
		Func:                   fn,
		AllowedUnauthenticated: false,
		AllowedUntrusted:       false,
		ValidPath:              false,
	}
}

// Unrestricted defines an unrestricted route.
func Unrestricted(fn Func) Route {
	return Route{
		Func:                   fn,
		AllowedUnauthenticated: true,
		AllowedUntrusted:       true,
		ValidPath:              false,
	}
}

// UntrustedValidPath defines a untrusted route.
func UntrustedValidPath(fn Func) Route {
	return Route{
		Func:                   fn,
		AllowedUnauthenticated: true,
		AllowedUntrusted:       true,
		ValidPath:              true,
	}
}
