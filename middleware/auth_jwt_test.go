package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aphistic/sweet"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type JWTAuthSuite struct{}

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

func (s *JWTAuthSuite) TestAuthorize(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.JSON(GetJWTClaims(ctx))
	}

	wrapped, err := NewAuthMiddleware(
		NewJWTAuthorizer(func(token *jwt.Token) (interface{}, error) {
			return jwtSigningSecret, nil
		}),
	).Convert(bare)

	Expect(err).To(BeNil())

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
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	_, body, _ := response.Serialize(resp)
	Expect(body).To(MatchJSON(fmt.Sprintf(`{
		"iss": "test",
		"foo": "bar",
		"exp": %d
	}`, inOneMinute)))
}

func (s *JWTAuthSuite) TestAuthorizeHeader(t sweet.T) {
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

	Expect(err).To(BeNil())

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
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	_, body, _ := response.Serialize(resp)
	Expect(body).To(MatchJSON(fmt.Sprintf(`{
		"iss": "test",
		"foo": "bar",
		"exp": %d
	}`, inOneMinute)))
}

func (s *JWTAuthSuite) TestExpiredToken(t sweet.T) {
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

	Expect(err).To(BeNil())

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
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("Token is expired"))
}

func (s *JWTAuthSuite) TestMalformedToken(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/test-auth?jwt=bad_jwt", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("token contains an invalid number of segments"))
}

func (s *JWTAuthSuite) TestNoToken(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("no token present in request"))
}

func (s *JWTAuthSuite) TestMismatchedSecret(t sweet.T) {
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

	Expect(err).To(BeNil())

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
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("signature is invalid"))
}

//
//

func makeJWTToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, _ := token.SignedString(jwtSigningSecret)
	return signedString
}
