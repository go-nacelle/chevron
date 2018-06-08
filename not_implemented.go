package chevron

import "context"

type tokenNotImplementedHandler string

var TokenNotImplementedHandler = tokenNotImplementedHandler("chevron.not_implemented_handler")

func GetNotImplementedHandler(ctx context.Context) Handler {
	if val, ok := ctx.Value(TokenNotImplementedHandler).(Handler); ok {
		return val
	}

	return defaultNotImplementedHandler
}

func setNotImplementedHandler(ctx context.Context, handler Handler) context.Context {
	return context.WithValue(ctx, TokenNotImplementedHandler, handler)
}
