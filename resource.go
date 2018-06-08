package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
)

type (
	Resource interface {
		Handle(context.Context, *http.Request, nacelle.Logger) response.Response
	}

	resource struct {
		hm     handlerMap
		router *router
	}
)

func (r *resource) Handle(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	if method, ok := methodMap[req.Method]; ok {
		if handler, ok := r.hm[method]; ok {
			return handler(ctx, req, logger)
		}
	}

	return r.router.notImplementedHandler(ctx, req, logger)
}
