package chevron

type (
	// MiddlewareConfig is a function that decorates a map from HTTP methods
	// to handlers.
	MiddlewareConfig func(handlerMap) error
)

// WithMiddleware applies the given middleware to all HTTP methods in the
// handler map.
func WithMiddleware(middleware Middleware) MiddlewareConfig {
	return func(hm handlerMap) error {
		return applyMiddleware(middleware, hm, allMethods)
	}
}

// WithMiddlewareFor applies the given middleware to the provided HTTP
// methods in the handler map.
func WithMiddlewareFor(middleware Middleware, methods ...Method) MiddlewareConfig {
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
