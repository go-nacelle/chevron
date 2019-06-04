package middleware

import (
	"context"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"

	"github.com/go-nacelle/chevron"
)

type (
	AuthMiddleware struct {
		authorizer                  Authorizer
		errorFactory                ErrorFactory
		forbiddenResponseFactory    ResponseFactory
		unauthorizedResponseFactory ErrorFactory
	}

	Authorizer interface {
		Authorize(context.Context, *http.Request) (AuthResult, interface{}, error)
	}

	AuthResult     int
	AuthorizerFunc func(context.Context, *http.Request) (AuthResult, interface{}, error)

	tokenAuthPayload string
)

var TokenAuthPayload = tokenAuthPayload("chevron.middleware.auth")

const (
	AuthResultInvalid AuthResult = iota
	AuthResultOK
	AuthResultForbidden
	AuthResultUnauthorized
)

func (f AuthorizerFunc) Authorize(ctx context.Context, req *http.Request) (AuthResult, interface{}, error) {
	return f(ctx, req)
}

func NewAuthMiddleware(authorizer Authorizer, configs ...AuthMiddlewareConfigFunc) chevron.Middleware {
	m := &AuthMiddleware{
		authorizer:                  authorizer,
		errorFactory:                defaultErrorFactory,
		forbiddenResponseFactory:    defaultForbiddenResponseFactory,
		unauthorizedResponseFactory: defaultUnauthorizedResponseFactory,
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *AuthMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
		result, payload, err := m.authorizer.Authorize(ctx, req)

		switch result {
		case AuthResultForbidden:
			return m.forbiddenResponseFactory()
		case AuthResultUnauthorized:
			return m.unauthorizedResponseFactory(err)
		default:
		}

		if err != nil {
			logger.Error("failed to invoke authorizer (%s)", err.Error())
			return m.errorFactory(err)
		}

		return f(context.WithValue(ctx, TokenAuthPayload, payload), req, logger)
	}

	return handler, nil
}

func defaultForbiddenResponseFactory() response.Response {
	return response.Empty(http.StatusForbidden)
}

func defaultUnauthorizedResponseFactory(err error) response.Response {
	return response.Empty(http.StatusUnauthorized)
}
