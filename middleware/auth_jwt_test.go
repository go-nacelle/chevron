package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

type CustomClaims struct {
	jwt.StandardClaims
	Foo string `json:"foo"`
}

var (
	jwtSigningSecret = []byte("super secret")
	now              = time.Now()
	inOneMinute      = now.Add(+time.Minute).Unix()
	oneMinuteAgo     = now.Add(-time.Minute).Unix()
)

func TestJWTAuthAuthorize(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(func(token *jwt.Token) (interface{}, error) {
			return jwtSigningSecret, nil
		}),
	).Convert(bare)

	assert.Nil(t, err)

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: inOneMinute,
			Issuer:    "test",
		},
		Foo: "bar",
	}

	url := fmt.Sprintf(
		"/test-auth?jwt=%s",
		makeJWTToken(claims),
	)

	r, _ := http.NewRequest("GET", url, nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	_, body, _ := response.Serialize(resp)

	expected := fmt.Sprintf(`{
		"iss": "test",
		"foo": "bar",
		"exp": %d
	}`, inOneMinute)
	assert.JSONEq(t, expected, string(body))
}

func TestJWTAuthAuthorizeHeader(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(
			func(token *jwt.Token) (interface{}, error) {
				return jwtSigningSecret, nil
			},
			WithJWTAuthExtractor(NewJWTHeaderExtractor("Authorization", "BEARER")),
		),
	).Convert(bare)

	assert.Nil(t, err)

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: inOneMinute,
			Issuer:    "test",
		},
		Foo: "bar",
	}

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	r.Header.Add("Authorization", fmt.Sprintf("BEARER %s", makeJWTToken(claims)))
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	_, body, _ := response.Serialize(resp)

	expected := fmt.Sprintf(`{
		"iss": "test",
		"foo": "bar",
		"exp": %d
	}`, inOneMinute)
	assert.JSONEq(t, expected, string(body))
}

func TestJWTAuthExpiredToken(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(
			func(token *jwt.Token) (interface{}, error) {
				return jwtSigningSecret, nil
			},
		),
		WithAuthUnauthorizedResponseFactory(func(err error) response.Response {
			return response.Respond([]byte(err.Error())).SetStatusCode(http.StatusUnauthorized)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: oneMinuteAgo,
			Issuer:    "test",
		},
		Foo: "bar",
	}

	url := fmt.Sprintf(
		"/test-auth?jwt=%s",
		makeJWTToken(claims),
	)

	r, _ := http.NewRequest("GET", url, nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "Token is expired", string(body))
}

func TestJWTAuthMalformedToken(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(
			func(token *jwt.Token) (interface{}, error) {
				return jwtSigningSecret, nil
			},
		),
		WithAuthUnauthorizedResponseFactory(func(err error) response.Response {
			return response.Respond([]byte(err.Error())).SetStatusCode(http.StatusUnauthorized)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/test-auth?jwt=bad_jwt", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "token contains an invalid number of segments", string(body))
}

func TestJWTAuthNoToken(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(
			func(token *jwt.Token) (interface{}, error) {
				return jwtSigningSecret, nil
			},
		),
		WithAuthUnauthorizedResponseFactory(func(err error) response.Response {
			return response.Respond([]byte(err.Error())).SetStatusCode(http.StatusUnauthorized)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "no token present in request", string(body))
}

func TestJWTAuthMismatchedSecret(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(func(token *jwt.Token) (interface{}, error) {
			return []byte("foobar"), nil
		}),
		WithAuthUnauthorizedResponseFactory(func(err error) response.Response {
			return response.Respond([]byte(err.Error())).SetStatusCode(http.StatusUnauthorized)
		}),
	).Convert(bare)

	assert.Nil(t, err)

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: inOneMinute,
			Issuer:    "test",
		},
		Foo: "bar",
	}

	url := fmt.Sprintf(
		"/test-auth?jwt=%s",
		makeJWTToken(claims),
	)

	r, _ := http.NewRequest("GET", url, nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())

	_, body, _ := response.Serialize(resp)
	assert.Equal(t, "signature is invalid", string(body))
}

//
//

func makeJWTToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, _ := token.SignedString(jwtSigningSecret)
	return signedString
}
