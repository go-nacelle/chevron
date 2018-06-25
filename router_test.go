package chevron

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type RouterSuite struct{}

func (s *RouterSuite) TestRegisterAddsToMux(t sweet.T) {
	var (
		container, _ = nacelle.MakeServiceContainer()
		router       = NewRouter(container)
	)

	// Register resource
	Expect(router.Register("/foo", &SimpleGetSpec{})).To(BeNil())

	// Matched route
	req, _ := http.NewRequest("GET", "/foo", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	Expect(recorder.Code).To(Equal(http.StatusNoContent))

	// Unmatched route
	req, _ = http.NewRequest("GET", "/bar", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	Expect(recorder.Code).To(Equal(http.StatusNotFound))
}

func (s *RouterSuite) TestRegisterInjectsServices(t sweet.T) {
	var (
		container, _ = nacelle.MakeServiceContainer()
		router       = NewRouter(container)
	)

	for _, val := range []string{"a", "b", "c"} {
		container.MustSet(val, val)
	}

	spec := &ServiceSpec{}
	Expect(router.Register("/users", spec)).To(BeNil())
	Expect(spec.A).To(Equal("a"))
	Expect(spec.B).To(Equal("b"))
	Expect(spec.C).To(Equal("c"))
}

func (s *RouterSuite) TestRegisterDuplicateURL(t sweet.T) {
	var (
		container, _ = nacelle.MakeServiceContainer()
		router       = NewRouter(container)
	)

	err1 := router.Register("/users", &EmptySpec{})
	err2 := router.Register("/users", &EmptySpec{})

	Expect(err1).To(BeNil())
	Expect(err2).NotTo(BeNil())
	Expect(err2.Error()).To(ContainSubstring("resource already registered"))
}

func (s *RouterSuite) TestRegisterWithMiddleware(t sweet.T) {
	var (
		container, _ = nacelle.MakeServiceContainer()
		router       = NewRouter(container)
		calls        = []string{}
	)

	middlewareFactory := func(name string) Middleware {
		return MiddlewareFunc(func(h Handler) (Handler, error) {
			handler := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
				calls = append(calls, name)
				return h(ctx, r, logger)
			}

			return handler, nil
		})
	}

	// Add global middleware
	router.AddMiddleware(middlewareFactory("a"))
	router.AddMiddleware(middlewareFactory("b"))
	router.AddMiddleware(middlewareFactory("c"))

	err := router.Register(
		"/foo",
		&SimpleGetSpec{},
		WithMiddleware(middlewareFactory("d")),
		WithMiddlewareFor(middlewareFactory("e"), MethodGet),
		WithMiddleware(middlewareFactory("f")),
	)

	Expect(err).To(BeNil())

	req, _ := http.NewRequest("GET", "/foo", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	Expect(recorder.Code).To(Equal(http.StatusNoContent))
	Expect(calls).To(Equal([]string{"a", "b", "c", "d", "e", "f"}))
}

//
//

type ServiceSpec struct {
	*EmptySpec

	A string `service:"a"`
	B string `service:"b"`
	C string `service:"c"`
}

type SimpleGetSpec struct {
	*EmptySpec
}

func (s *SimpleGetSpec) Get(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
	return response.Empty(http.StatusNoContent)
}
