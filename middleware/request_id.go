package middleware

import (
	"context"
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	"github.com/google/uuid"

	"github.com/efritz/chevron"
)

type (
	RequestIDMiddleware struct {
		requestIDGenerator RequestIDGenerator
		errorFactory       ErrorFactory
	}

	RequestIDGenerator func() (string, error)
	tokenRequestID     string
)

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
func NewRequestID(configs ...RequestIDConfigFunc) *RequestIDMiddleware {
	m := &RequestIDMiddleware{
		requestIDGenerator: defaultRequestIDGenerator,
		errorFactory:       defaultErrorFactory,
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *RequestIDMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
		requestID, err := m.getIDFromRequest(req)
		if err != nil {
			logger.Error("Failed to generate request ID (%s)", err.Error())
			return m.errorFactory(err)
		}

		wrappedCtx := context.WithValue(
			ctx,
			TokenRequestID,
			requestID,
		)

		wrappedLogger := logger.WithFields(nacelle.LogFields{
			"request_id": requestID,
		})

		resp := f(wrappedCtx, req, wrappedLogger)
		resp.SetHeader("X-Request-ID", requestID)
		return resp
	}

	return handler, nil
}

func (m *RequestIDMiddleware) getIDFromRequest(req *http.Request) (string, error) {
	if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
		return requestID, nil
	}

	return m.requestIDGenerator()
}

func defaultRequestIDGenerator() (string, error) {
	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return raw.String(), nil
}
