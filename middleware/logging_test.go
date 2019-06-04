package middleware

//go:generate go-mockgen github.com/go-nacelle/nacelle -i Logger -d mocks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	. "github.com/onsi/gomega"

	"github.com/go-nacelle/chevron/middleware/mocks"
)

type LoggingSuite struct{}

func (s *LoggingSuite) TestHandle(t sweet.T) {
	var (
		clock           = glock.NewMockClock()
		logger          = mocks.NewMockLogger()
		decoratedLogger = mocks.NewMockLogger()
	)

	logger.WithFieldsFunc = func(nacelle.LogFields) nacelle.Logger {
		return decoratedLogger
	}

	handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		Expect(logger).To(Equal(decoratedLogger))
		clock.Advance(time.Millisecond * 322)
		return response.Empty(http.StatusOK)
	}

	wrapped, err := NewLogging(WithLoggingClock(clock)).Convert(handler)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	r.RequestURI = "/foo/bar"
	r.RemoteAddr = "250.34.143.226"
	r.Header.Add("User-Agent", "chevron-test")

	resp := wrapped(context.Background(), r, logger)
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	// Check function calls
	Expect(logger.WithFieldsFuncCallCount()).To(Equal(1))
	Expect(decoratedLogger.InfoFuncCallCount()).To(Equal(2))

	// Check request logger fields
	Expect(logger.WithFieldsFuncCallParams()[0].Arg0).To(Equal(nacelle.LogFields{
		"method":       "GET",
		"request_uri":  "/foo/bar",
		"client":       "250.34.143.226",
		"request_host": "example.com",
		"user_agent":   "chevron-test",
	}))

	// Check access log messages
	params1 := decoratedLogger.InfoFuncCallParams()[0]
	params2 := decoratedLogger.InfoFuncCallParams()[1]
	Expect(fmt.Sprintf(params1.Arg0, params1.Arg1...)).To(Equal("Handling HTTP request GET /foo/bar"))
	Expect(fmt.Sprintf(params2.Arg0, params2.Arg1...)).To(Equal("Handled HTTP request GET /foo/bar -> 200 in 322ms"))
}
