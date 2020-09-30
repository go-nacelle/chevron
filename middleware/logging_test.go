package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/derision-test/glock"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"

	"github.com/go-nacelle/chevron/middleware/mocks"
)

func TestLoggingHandle(t *testing.T) {
	var (
		clock           = glock.NewMockClock()
		logger          = mocks.NewMockLogger()
		decoratedLogger = mocks.NewMockLogger()
	)

	logger.WithFieldsFunc.SetDefaultHook(func(nacelle.LogFields) nacelle.Logger {
		return decoratedLogger
	})

	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		assert.Equal(t, decoratedLogger, logger)
		clock.Advance(time.Millisecond * 322)
		return response.Empty(http.StatusOK)
	}

	wrapped, err := NewLogging(WithLoggingClock(clock)).Convert(handler)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	r.RequestURI = "/foo/bar"
	r.RemoteAddr = "250.34.143.226"
	r.Header.Add("User-Agent", "chevron-test")

	resp := wrapped(context.Background(), r, logger)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check function calls
	assert.Equal(t, 1, len(logger.WithFieldsFunc.History()))
	assert.Equal(t, 2, len(decoratedLogger.InfoFunc.History()))

	expected := nacelle.LogFields{
		"method":       "GET",
		"request_uri":  "/foo/bar",
		"client":       "250.34.143.226",
		"request_host": "example.com",
		"user_agent":   "chevron-test",
	}
	assert.Equal(t, expected, logger.WithFieldsFunc.History()[0].Arg0)

	// Check access log messages
	params1 := decoratedLogger.InfoFunc.History()[0]
	params2 := decoratedLogger.InfoFunc.History()[1]
	assert.Equal(t, "Handling HTTP request GET /foo/bar", fmt.Sprintf(params1.Arg0, params1.Arg1...))
	assert.Equal(t, "Handled HTTP request GET /foo/bar -> 200 in 322ms", fmt.Sprintf(params2.Arg0, params2.Arg1...))
}
