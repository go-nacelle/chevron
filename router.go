package chevron

import (
	"context"
	"fmt"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/gorilla/mux"
)

// Router stores
type Router interface {
	http.Handler

	// AddMiddleware registers middleware for all resources registered after
	// the invocation of this method.
	AddMiddleware(middleware Middleware)

	// Register creates a resource from the given resource spec and set of
	// middleware instances and registers it to the given URL pattern.
	Register(url string, spec ResourceSpec, configs ...MiddlewareConfigFunc) error

	// MustRegister calls Register and panics on error.
	MustRegister(url string, spec ResourceSpec, configs ...MiddlewareConfigFunc)

	RegisterHandler(url string, handler http.Handler)
}

type router struct {
	services              nacelle.ServiceContainer
	logger                nacelle.Logger
	middleware            []Middleware
	mux                   *mux.Router
	resources             map[string]struct{}
	notFoundHandler       Handler
	notImplementedHandler Handler
	baseCtx               context.Context
}

type handlerMap map[Method]Handler

// NewRouter creates a new router.
func NewRouter(services nacelle.ServiceContainer, logger nacelle.Logger, configs ...RouterConfigFunc) Router {
	r := &router{
		services:              services,
		logger:                logger,
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

// AddMiddleware registers middleware for all resources. This middleware is not
// retroactively applied to resources which have already been registered (so
// attention to the order in which middleware is registered is required on the
// part of the developer).
func (r *router) AddMiddleware(middleware Middleware) {
	r.middleware = append(r.middleware, middleware)
}

// Register creates a resource from the given resource spec and set of
// middleware instances and registers it to the given URL pattern. It
// is an error to register the same URL pattern twice.
func (r *router) Register(url string, spec ResourceSpec, configs ...MiddlewareConfigFunc) error {
	if _, ok := r.resources[url]; ok {
		return fmt.Errorf("resource already registered to url pattern `%s`", url)
	}

	r.resources[url] = struct{}{}

	if err := r.services.Inject(spec); err != nil {
		return err
	}

	resource, err := r.decorateResource(spec, configs...)
	if err != nil {
		return err
	}

	r.RegisterHandler(url, convert(r.baseCtx, resource.Handle, r.logger))
	return nil
}

func (r *router) decorateResource(spec ResourceSpec, configs ...MiddlewareConfigFunc) (Resource, error) {
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

// MustRegister calls Register and panics on error.
func (r *router) MustRegister(url string, spec ResourceSpec, configs ...MiddlewareConfigFunc) {
	if err := r.Register(url, spec, configs...); err != nil {
		panic(err.Error())
	}
}

func (r *router) RegisterHandler(url string, handler http.Handler) {
	r.mux.Handle(url, handler)
}

// ServeHTTP invokes the handler registered to the request URL and
// writes the response to the given ResponseWriter.
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
