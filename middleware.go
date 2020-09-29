package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
)

// Handler converts an HTTP request into a response object. The handler
// also has access to a request context populated by the router and
// registered middleware as well as a request logger which may also be
// decorated by registered middleware.
type Handler func(context.Context, *http.Request, nacelle.Logger) response.Response

// Middleware transforms a handler into another decorated handler.
type Middleware interface {
	// Convert applies the middleware transformation to a handler. If
	// the middleware is supplied invalid arguments it may return an
	// error at the time of decoration.
	Convert(Handler) (Handler, error)
}

// MiddlewareFunc is signature for single-function middleware.
type MiddlewareFunc func(Handler) (Handler, error)

func (f MiddlewareFunc) Convert(h Handler) (Handler, error) {
	return f(h)
}
