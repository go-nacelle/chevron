package chevron

// MiddlewareConfigFunc is a function that decorates a map from HTTP methods
// to handlers.
type MiddlewareConfigFunc func(handlerMap) error

// WithMiddleware applies the given middleware to all HTTP methods in the
// handler map.
func WithMiddleware(middleware Middleware) MiddlewareConfigFunc {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, allMethods)
	}
}

// WithMiddlewareFor applies the given middleware to the provided HTTP
// methods in the handler map.
func WithMiddlewareFor(middleware Middleware, methods ...Method) MiddlewareConfigFunc {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, methods)
	}
}

func applyMiddleware(middleware Middleware, hm handlerMap, methods []Method) error {
	for _, method := range methods {
		wrapped, err := middleware.Convert(hm[method])
		if err != nil {
			return err
		}

		hm[method] = wrapped
	}

	return nil
}
