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

// TokenRequestID is the unique token to which the current request's unique
// ID is written to the request context.
var TokenRequestID = tokenRequestID("chevron.middleware.request_id")

// GetRequestID retrieves the current request's unique ID. If no request ID
// is registered with this context, the empty string is returned.
func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(TokenRequestID).(string); ok {
		return val
	}

	return ""
}

// NewRequestID creates middleware that generates a unique ID for the request.
// If the header X-Request-ID is present in the request, that value is used
// instead. The request ID is added to the context and to logger attributes,
// and the X-Request-ID header is added to the wrapped handler's resulting
// response.
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
