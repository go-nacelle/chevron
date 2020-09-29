package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
)

// Resource represents an object or behavior accessible from a unique URL
// pattern.
type Resource interface {
	// Handle is invoked when the resource is requested, regardless of the
	// HTTP method.
	Handle(context.Context, *http.Request, nacelle.Logger) response.Response
}

type resource struct {
	hm     handlerMap
	router *router
}

// Handle invokes the correct handler based on HTTP method, or the router's not
// implemented handler if no handler for that method is registered.
func (r *resource) Handle(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	if method, ok := methodMap[req.Method]; ok {
		if handler, ok := r.hm[method]; ok {
			return handler(ctx, req, logger)
		}
	}

	return r.router.notImplementedHandler(ctx, req, logger)
}
