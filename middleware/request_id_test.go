package middleware

import (
	"context"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type RequestIDSuite struct{}

func (s *RequestIDSuite) TestGetRequestID(t sweet.T) {
	var (
		val = "x"
		ctx = context.WithValue(context.Background(), TokenRequestID, val)
	)

	// Present
	Expect(GetRequestID(ctx)).To(Equal(val))

	// Missing
	Expect(GetRequestID(context.Background())).To(BeEmpty())
}

func (s *RequestIDSuite) TestNewRequestIDGenerated(t sweet.T) {
	var ctxVal string
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetRequestID(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewRequestID().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))
	Expect(resp.Header("X-Request-ID")).To(Equal(ctxVal))
	Expect(resp.Header("X-Request-ID")).To(HaveLen(36))
}

func (s *RequestIDSuite) TestNewRequestIDSuppliedByClient(t sweet.T) {
	var ctxVal string
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetRequestID(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewRequestID().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("X-Request-ID", "1234")
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))
	Expect(resp.Header("X-Request-ID")).To(Equal(ctxVal))
	Expect(resp.Header("X-Request-ID")).To(Equal("1234"))
}
