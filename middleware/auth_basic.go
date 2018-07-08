package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/efritz/response"
)

type (
	basicAuthorizer struct {
		validator BasicAuthValidator
	}

	BasicAuthValidator func(context.Context, string, string) (bool, error)
)

func GetBasicAuthUsername(ctx context.Context) string {
	if val, ok := ctx.Value(TokenAuthPayload).(string); ok {
		return val
	}

	return ""
}

func NewBasicAuthorizer(validator BasicAuthValidator) Authorizer {
	return &basicAuthorizer{
		validator: validator,
	}
}

func (a *basicAuthorizer) Authorize(ctx context.Context, req *http.Request) (AuthResult, interface{}, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return AuthResultUnauthorized, nil, nil
	}

	auth, err := a.validator(ctx, username, password)
	if err != nil {
		return AuthResultInvalid, nil, err
	}

	if !auth {
		return AuthResultForbidden, nil, nil
	}

	return AuthResultOK, username, nil
}

func defaultBasicAuthValidator(ctx context.Context, username, password string) (bool, error) {
	return false, nil
}

func NewBasicUnauthorizedResponseFactory(realm string) ErrorFactory {
	return func(err error) response.Response {
		resp := response.Empty(http.StatusUnauthorized)
		resp.SetHeader("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
		return resp
	}
}
