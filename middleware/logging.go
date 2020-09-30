package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/derision-test/glock"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"

	"github.com/go-nacelle/chevron"
)

type LoggingMiddleware struct {
	clock glock.Clock
}

func NewLogging(configs ...LoggingConfigFunc) chevron.Middleware {
	m := &LoggingMiddleware{
		clock: glock.NewRealClock(),
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *LoggingMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
		logger = logger.WithFields(nacelle.LogFields{
			"method":       req.Method,
			"request_uri":  req.RequestURI,
			"client":       req.RemoteAddr,
			"request_host": req.Host,
			"user_agent":   req.UserAgent(),
		})

		logger.Info(
			"Handling HTTP request %s %s",
			req.Method,
			req.RequestURI,
		)

		start := m.clock.Now()
		resp := f(ctx, req, logger)
		duration := int(m.clock.Now().Sub(start) / time.Millisecond)

		logger.Info(
			"Handled HTTP request %s %s -> %d in %dms",
			req.Method,
			req.RequestURI,
			resp.StatusCode(),
			duration,
		)

		return resp
	}

	return handler, nil
}
