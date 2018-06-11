package chevron

import (
	"github.com/efritz/nacelle"
)

// RouterConfig is a function used to initialize a new router.
type RouterConfig func(*router)

// WithLogger sets the router's logger.
func WithLogger(logger nacelle.Logger) RouterConfig {
	return func(r *router) { r.logger = logger }
}

// WithNotFoundHandler sets the handler invoked when a requested
// URL cannot be matched with any registered URL pattern.
func WithNotFoundHandler(handler Handler) RouterConfig {
	return func(r *router) { r.notFoundHandler = handler }
}

// WithNotImplementedHandler sets the handler invoked when a
// resource does not implemented the requested HTTP verb.
func WithNotImplementedHandler(handler Handler) RouterConfig {
	return func(r *router) { r.notImplementedHandler = handler }
}
