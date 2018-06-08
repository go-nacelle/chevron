package chevron

import (
	"github.com/efritz/nacelle"
)

type RouterConfig func(*router)

func WithLogger(logger nacelle.Logger) RouterConfig {
	return func(r *router) { r.logger = logger }
}

func WithNotFoundHandler(handler Handler) RouterConfig {
	return func(r *router) { r.notFoundHandler = handler }
}

func WithNotImplementedHandler(handler Handler) RouterConfig {
	return func(r *router) { r.notImplementedHandler = handler }
}
