package middleware

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"

	"github.com/efritz/chevron"
)

// NewRecovery creates middleware that captures panics from the handler
// and converts them to 500-level responses. The value of the panic is
// logged at error level.
func NewRecovery() chevron.Middleware {
	return chevron.MiddlewareFunc(func(f chevron.Handler) (chevron.Handler, error) {
		handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) (resp response.Response) {
			defer func() {
				if err := recover(); err != nil {
					// TODO - should add stack in logs
					logger.Error("Request handler panicked (%s)", err)
					resp = response.Empty(http.StatusInternalServerError)
				}
			}()

			resp = f(ctx, req, logger)
			return
		}

		return handler, nil
	})
}
