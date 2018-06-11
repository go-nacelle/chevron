package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
)

type (
	// ResourceSpec represents the set of handlers (one for each HTTP method) to
	// which a resource can respond.
	ResourceSpec interface {
		// Get is the handler invoked for the GET HTTP method.
		Get(context.Context, *http.Request, nacelle.Logger) response.Response

		// Options is the handler invoked for the OPTIONS HTTP method.
		Options(context.Context, *http.Request, nacelle.Logger) response.Response

		// Post is the handler invoked for the POST HTTP method.
		Post(context.Context, *http.Request, nacelle.Logger) response.Response

		// Put is the handler invoked for the PUT HTTP method.
		Put(context.Context, *http.Request, nacelle.Logger) response.Response

		// Patch is the handler invoked for the PATCH HTTP method.
		Patch(context.Context, *http.Request, nacelle.Logger) response.Response

		// Delete is the handler invoked for the DELETE HTTP method.
		Delete(context.Context, *http.Request, nacelle.Logger) response.Response
	}

	// EmptySpec is a complete implementation of ResourceSpec that invokes the
	// router's not implemented handler. A pointer to this struct should be the
	// first embedded field in any resource - this allows a resource to simply
	// "override" the handlers for methods relevant to a resource.
	EmptySpec struct{}
)

// Get invokes the router's not implemented handler.
func (es *EmptySpec) Get(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

// Options invokes the router's not implemented handler.
func (es *EmptySpec) Options(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

// Post invokes the router's not implemented handler.
func (es *EmptySpec) Post(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

// Put invokes the router's not implemented handler.
func (es *EmptySpec) Put(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

// Patch invokes the router's not implemented handler.
func (es *EmptySpec) Patch(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

// Delete invokes the router's not implemented handler.
func (es *EmptySpec) Delete(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}
