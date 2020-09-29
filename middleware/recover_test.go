package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/chevron/middleware/mocks"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

func TestRecoverBaseline(t *testing.T) {
	// This test ensures that a handler that panics does not have
	// the same behavior with this middleware enabled. Here we just
	// show the default behavior.

	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		panic("oops")
	}

	r, _ := http.NewRequest("GET", "/", nil)
	assert.Panics(t, func() { handler(context.Background(), r, nacelle.NewNilLogger()) })
}

func TestRecoverWithRecover(t *testing.T) {
	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		var w io.Writer
		w.Write(nil)
		return nil
	}

	wrapped, err := NewRecovery().Convert(handler)
	assert.Nil(t, err)

	logger := mocks.NewMockLogger()

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, logger)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	assert.Equal(t, 1, len(logger.ErrorFunc.History()))

	var (
		args       = logger.ErrorFunc.History()[0].Arg1
		message    = fmt.Sprintf("%s", args[0])
		stacktrace = fmt.Sprintf("%s", args[1])
	)

	// Note: line number may change if code is added above of this test
	assert.Equal(t, "runtime error: invalid memory address or nil pointer dereference", message)
	assert.Contains(t, stacktrace, "go-nacelle/chevron/middleware/recover_test.go:32")
}

func TestRecoverWithRecoverErrorFactory(t *testing.T) {
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

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, "whoopsie", arg)
	assert.Equal(t, expectedResp, resp)
}
