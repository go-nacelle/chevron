package chevron

import (
	"context"
	"net/http"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&RouterSuite{})
		s.AddSuite(&MiddlewareOptionsSuite{})
		s.AddSuite(&SpecSuite{})
		s.AddSuite(&ResourceSuite{})
	})
}

//
//

func makeEmptyHandler(status int) Handler {
	return func(context.Context, *http.Request, nacelle.Logger) response.Response {
		return response.Empty(status)
	}
}
