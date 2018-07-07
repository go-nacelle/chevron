package middleware

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&CacheSuite{})
		s.AddSuite(&GzipSuite{})
		s.AddSuite(&LoggingSuite{})
		s.AddSuite(&RecoverSuite{})
		s.AddSuite(&RequestIDSuite{})
		s.AddSuite(&SchemaSuite{})
	})
}
