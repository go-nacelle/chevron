package middleware

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"

	"github.com/efritz/chevron"
)

func NewRecovery() chevron.Middleware {
	return func(f chevron.Handler) (chevron.Handler, error) {
		handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) (resp response.Response) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Request handler panic'd (%s)", err)
					resp = response.Empty(http.StatusInternalServerError)
				}
			}()

			resp = f(ctx, req, logger)
			return
		}

		return handler, nil
	}
}
