package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

var testBasicValidator = func(ctx context.Context, u, p string) (bool, error) {
	return u == "admin" && p == "secret", nil
}

func TestBasicAuthAuthorize(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.Respond([]byte(GetBasicAuthUsername(ctx)))
	}

	wrapped, err := NewAuthMiddleware(NewBasicAuthorizer(testBasicValidator)).Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	r.SetBasicAuth("admin", "secret")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "admin", string(body))
}

func TestBasicAuthAuthorizeBadAuth(t *testing.T) {
	called := false

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusOK)
	}

	wrapped, err := NewAuthMiddleware(
		NewBasicAuthorizer(testBasicValidator),
		WithAuthForbiddenResponseFactory(func() response.Response {
			return response.Respond([]byte("403")).SetStatusCode(http.StatusForbidden)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	r.SetBasicAuth("admin", "old-secret")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusForbidden, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "403", string(body))
	assert.False(t, called)
}

func TestBasicAuthAuthorizeMissingAuth(t *testing.T) {
	called := false

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusOK)
	}

	wrapped, err := NewAuthMiddleware(
		NewBasicAuthorizer(testBasicValidator),
		WithAuthUnauthorizedResponseFactory(func(err error) response.Response {
			return response.Respond([]byte("401")).SetStatusCode(http.StatusUnauthorized)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/test-auth", nil)

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "401", string(body))
	assert.False(t, called)
}

func TestBasicAuthDefaultUnauthorizedResponseFactory(t *testing.T) {
	resp := NewBasicUnauthorizedResponseFactory("test")(fmt.Errorf("utoh"))
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	assert.Equal(t, `Basic realm="test"`, resp.Header("WWW-Authenticate"))
}
