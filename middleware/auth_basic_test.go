package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type BasicAuthSuite struct{}

var testBasicValidator = func(ctx context.Context, u, p string) (bool, error) {
	return u == "admin" && p == "secret", nil
}

func (s *BasicAuthSuite) TestAuthorize(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.Respond([]byte(GetBasicAuthUsername(ctx)))
	}

	wrapped, err := NewAuthMiddleware(NewBasicAuthorizer(testBasicValidator)).Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	r.SetBasicAuth("admin", "secret")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("admin"))
}

func (s *BasicAuthSuite) TestAuthorizeBadAuth(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/test-auth", nil)
	r.SetBasicAuth("admin", "old-secret")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusForbidden))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("403"))
	Expect(called).To(BeFalse())
}

func (s *BasicAuthSuite) TestAuthorizeMissingAuth(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/test-auth", nil)

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))

	_, body, _ := response.Serialize(resp)
	Expect(string(body)).To(Equal("401"))
	Expect(called).To(BeFalse())
}

func (s *BasicAuthSuite) TestDefaultUnauthorizedResponseFactory(t sweet.T) {
	resp := NewBasicUnauthorizedResponseFactory("test")(fmt.Errorf("utoh"))
	Expect(resp.StatusCode()).To(Equal(http.StatusUnauthorized))
	Expect(resp.Header("WWW-Authenticate")).To(Equal(`Basic realm="test"`))
}
