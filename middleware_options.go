package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
)

type (
	Handler    func(context.Context, *http.Request, nacelle.Logger) response.Response
	Middleware func(Handler) (Handler, error)
)

func WithMiddleware(middleware Middleware) MiddlewareConfig {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, allMethods)
	}
}

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
