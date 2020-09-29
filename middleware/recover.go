package middleware

import (
	"context"
	"net/http"
	"runtime"

	"github.com/efritz/response"
	"github.com/go-nacelle/chevron"
	"github.com/go-nacelle/nacelle"
)

type RecoverMiddleware struct {
	errorFactory     PanicErrorFactory
	stackBufferSize  int
	logAllGoroutines bool
}

// NewRecovery creates middleware that captures panics from the handler
// and converts them to 500-level responses. The value of the panic is
// logged at error level.
func NewRecovery(configs ...RecoverConfigFunc) chevron.Middleware {
	m := &RecoverMiddleware{
		errorFactory:     defaultPanicErrorFactory,
		stackBufferSize:  4 << 10,
		logAllGoroutines: false,
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *RecoverMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) (resp response.Response) {
		defer func() {
			if err := recover(); err != nil {
				var (
					stack  = make([]byte, m.stackBufferSize)
					length = runtime.Stack(stack, m.logAllGoroutines)
				)

				logger.Error(
					"Request handler panicked (%s):\n%s",
					err,
					stack[:length],
				)

				resp = m.errorFactory(err)
			}
		}()

		resp = f(ctx, req, logger)
		return
	}

	return handler, nil
}
