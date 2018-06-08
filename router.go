package chevron

import (
	"context"
	"fmt"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/log"
	"github.com/efritz/response"
	"github.com/gorilla/mux"
)

type (
	Router interface {
		http.Handler

		AddMiddleware(middleware Middleware)
		Register(url string, spec ResourceSpec, configs ...MiddlewareConfig) error
		MustRegister(url string, spec ResourceSpec, configs ...MiddlewareConfig)
	}

	router struct {
		container             *nacelle.ServiceContainer
		logger                nacelle.Logger
		middleware            []Middleware
		mux                   *mux.Router
		resources             map[string]struct{}
		notFoundHandler       Handler
		notImplementedHandler Handler
		baseCtx               context.Context
	}

	handlerMap       map[Method]Handler
	MiddlewareConfig func(handlerMap) error
)

func NewRouter(container *nacelle.ServiceContainer, configs ...RouterConfig) Router {
	r := &router{
		container:             container,
		logger:                log.NewNilLogger(),
		mux:                   mux.NewRouter(),
		resources:             map[string]struct{}{},
		notFoundHandler:       defaultNotFoundHandler,
		notImplementedHandler: defaultNotImplementedHandler,
	}

	for _, config := range configs {
		config(r)
	}

	r.baseCtx = setNotImplementedHandler(context.Background(), r.notImplementedHandler)
	r.mux.NotFoundHandler = convert(r.baseCtx, r.notFoundHandler, r.logger)
	return r
}

func (r *router) AddMiddleware(middleware Middleware) {
	r.middleware = append(r.middleware, middleware)
}

func (r *router) Register(url string, spec ResourceSpec, configs ...MiddlewareConfig) error {
	if _, ok := r.resources[url]; ok {
		return fmt.Errorf("resource already registered to url pattern `%s`", url)
	}

	r.resources[url] = struct{}{}

	if err := r.container.Inject(spec); err != nil {
		return err
	}

	resource, err := r.decorateResource(spec, configs...)
	if err != nil {
		return err
	}

	r.mux.Handle(url, convert(r.baseCtx, resource.Handle, r.logger))
	return nil
}

func (r *router) decorateResource(spec ResourceSpec, configs ...MiddlewareConfig) (Resource, error) {
	hm := handlerMap{
		MethodGet:     spec.Get,
		MethodOptions: spec.Options,
		MethodPost:    spec.Post,
		MethodPut:     spec.Put,
		MethodPatch:   spec.Patch,
		MethodDelete:  spec.Delete,
	}

	for i := len(configs) - 1; i >= 0; i-- {
		if err := configs[i](hm); err != nil {
			return nil, err
		}
	}

	for i := len(r.middleware) - 1; i >= 0; i-- {
		if err := applyMiddleware(r.middleware[i], hm, allMethods); err != nil {
			return nil, err
		}
	}

	return &resource{hm: hm, router: r}, nil
}

func (r *router) MustRegister(url string, spec ResourceSpec, configs ...MiddlewareConfig) {
	if err := r.Register(url, spec, configs...); err != nil {
		panic(err.Error())
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

//
//

func convert(ctx context.Context, handler Handler, logger nacelle.Logger) http.Handler {
	responseHandler := func(r *http.Request) response.Response {
		return handler(ctx, r, logger)
	}

	return http.HandlerFunc(response.Convert(responseHandler))
}

func defaultNotFoundHandler(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
	return response.Empty(http.StatusNotFound)
}

func defaultNotImplementedHandler(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
	return response.Empty(http.StatusMethodNotAllowed)
}
