package middleware

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/log"
	"github.com/efritz/response"
	"github.com/satori/go.uuid"

	"github.com/efritz/chevron"
)

type tokenRequestID string

var TokenRequestID = tokenRequestID("chevron.middleware.request_id")

func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(TokenRequestID).(string); ok {
		return val
	}

	return ""
}

func NewRequestID() chevron.Middleware {
	return func(f chevron.Handler) (chevron.Handler, error) {
		handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
			requestID, err := getRequestIDFromRequest(req)
			if err != nil {
				// TODO
				return response.Empty(http.StatusInternalServerError)
			}

			wrappedCtx := context.WithValue(
				ctx,
				TokenRequestID,
				requestID,
			)

			wrappedLogger := logger.WithFields(log.Fields{
				"request_id": requestID,
			})

			resp := f(wrappedCtx, req, wrappedLogger)
			resp.SetHeader("X-Request-ID", requestID)
			return resp
		}

		return handler, nil
	}
}

func getRequestIDFromRequest(req *http.Request) (string, error) {
	if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
		return requestID, nil
	}

	raw, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return raw.String(), nil
}
