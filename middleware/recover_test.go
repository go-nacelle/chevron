package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/chevron/middleware/mocks"
	"github.com/efritz/nacelle"
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
	Expect(func() { handler(context.Background(), r, nacelle.NewNilLogger()) }).To(Panic())
}

func (s *RecoverSuite) TestWithRecover(t sweet.T) {
	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		var w io.Writer
		w.Write(nil)
		return nil
	}

	wrapped, err := NewRecovery().Convert(handler)
	Expect(err).To(BeNil())

	logger := mocks.NewMockLogger()

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, logger)
	Expect(resp.StatusCode()).To(Equal(http.StatusInternalServerError))
	Expect(logger.ErrorFuncCallCount()).To(Equal(1))

	var (
		args       = logger.ErrorFuncCallParams()[0].Arg1
		message    = fmt.Sprintf("%s", args[0])
		stacktrace = fmt.Sprintf("%s", args[1])
	)

	// Note: line number may change if code is added above of this test
	Expect(message).To(Equal("runtime error: invalid memory address or nil pointer dereference"))
	Expect(stacktrace).To(ContainSubstring("github.com/efritz/chevron/middleware/recover_test.go:34"))
}

func (s *RecoverSuite) TestWithRecoverErrorFactory(t sweet.T) {
	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		panic("whoopsie")
	}

	expectedResp := response.JSON(map[string]string{
		"message": "handler whoopsie",
	})

	var arg interface{}
	errorFactory := func(val interface{}) response.Response {
		arg = val
		return expectedResp
	}

	wrapped, err := NewRecovery(
		WithRecoverErrorFactory(errorFactory),
	).Convert(handler)

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(arg).To(Equal("whoopsie"))
	Expect(resp).To(Equal(expectedResp))
}
