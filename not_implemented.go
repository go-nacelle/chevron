package chevron

import "context"

type tokenNotImplementedHandler string

// TokenNotImplementedHandler is the unique token to which the router's
// not implemented handler is written to the request context.
var TokenNotImplementedHandler = tokenNotImplementedHandler("chevron.not_implemented_handler")

// GetNotImplementedHandler retrieves the router's not implemented handler
// from the given context. If no handler is registered with this context,
// the default handler is returned.
func GetNotImplementedHandler(ctx context.Context) Handler {
	if val, ok := ctx.Value(TokenNotImplementedHandler).(Handler); ok {
		return val
	}

	return defaultNotImplementedHandler
}

func setNotImplementedHandler(ctx context.Context, handler Handler) context.Context {
	return context.WithValue(ctx, TokenNotImplementedHandler, handler)
}
