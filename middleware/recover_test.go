package middleware

import (
	"context"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/log"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type RecoverSuite struct{}

func (s *RecoverSuite) TestBaseline(t sweet.T) {
	// This test ensures that a handler that panics does not have
	// the same behavior with this middleware enabled. Here we just
	// show the default behavior.

	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		panic("oops")
	}

	r, _ := http.NewRequest("GET", "/", nil)
	Expect(func() { handler(context.Background(), r, log.NewNilLogger()) }).To(Panic())
}

func (s *RecoverSuite) TestWithRecover(t sweet.T) {
	handler, err := NewRecovery()(func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		panic("oops")
	})

	Expect(err).To(BeNil())
	r, _ := http.NewRequest("GET", "/", nil)
	resp := handler(context.Background(), r, log.NewNilLogger())
	Expect(resp.StatusCode()).To(Equal(http.StatusInternalServerError))
}