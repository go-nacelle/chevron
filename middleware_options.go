package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
)

type (
	// Handler converts an HTTP request into a response object. The handler
	// also has access to a request context populated by the router and
	// registered middleware as well as a request logger which may also be
	// decorated by registered middleware.
	Handler func(context.Context, *http.Request, nacelle.Logger) response.Response

	// Middleware transforms a handler into another decorated handler. If
	// the middleware is supplied invalid arguments it may return an error
	// at the time of decoration.
	Middleware func(Handler) (Handler, error)

	// MiddlewareConfig is a function that decorates a map from HTTP methods
	// to handlers.
	MiddlewareConfig func(handlerMap) error
)

// WithMiddleware applies the given middleware to all HTTP methods in the
// handler map.
func WithMiddleware(middleware Middleware) MiddlewareConfig {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, allMethods)
	}
}

// WithMiddlewareFor applies the given middleware to the provided HTTP
// methods in the handler map.
func WithMiddlewareFor(middleware Middleware, methods ...Method) MiddlewareConfig {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, methods)
	}
}

func applyMiddleware(middleware Middleware, hm handlerMap, methods []Method) error {
	for _, method := range methods {
		wrapped, err := middleware(hm[method])
		if err != nil {
			return err
		}

		hm[method] = wrapped
	}

	return nil
}
