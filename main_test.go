package chevron

import (
	"context"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
)

func makeEmptyHandler(status int) Handler {
	return func(context.Context, *http.Request, nacelle.Logger) response.Response {
		return response.Empty(status)
	}
}
