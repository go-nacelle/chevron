package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
)

type (
	ResourceSpec interface {
		Get(context.Context, *http.Request, nacelle.Logger) response.Response
		Options(context.Context, *http.Request, nacelle.Logger) response.Response
		Post(context.Context, *http.Request, nacelle.Logger) response.Response
		Put(context.Context, *http.Request, nacelle.Logger) response.Response
		Patch(context.Context, *http.Request, nacelle.Logger) response.Response
		Delete(context.Context, *http.Request, nacelle.Logger) response.Response
	}

	EmptySpec struct{}
)

func (es *EmptySpec) Get(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

func (es *EmptySpec) Options(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

func (es *EmptySpec) Post(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

func (es *EmptySpec) Put(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

func (es *EmptySpec) Patch(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}

func (es *EmptySpec) Delete(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return GetNotImplementedHandler(ctx)(ctx, req, logger)
}
