package chevron

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

func TestRouterRegisterAddsToMux(t *testing.T) {
	var (
		container = nacelle.NewServiceContainer()
		logger    = nacelle.NewNilLogger()
		router    = NewRouter(container, logger)
	)

	// Register resource
	assert.Nil(t, router.Register("/foo", &SimpleGetSpec{}))

	// Matched route
	req, _ := http.NewRequest("GET", "/foo", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNoContent, recorder.Code)

	// Unmatched route
	req, _ = http.NewRequest("GET", "/bar", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestRouterRegisterInjectsServices(t *testing.T) {
	var (
		container = nacelle.NewServiceContainer()
		logger    = nacelle.NewNilLogger()
		router    = NewRouter(container, logger)
	)

	for _, val := range []string{"a", "b", "c"} {
		container.MustSet(val, val)
	}

	spec := &ServiceSpec{}
	assert.Nil(t, router.Register("/users", spec))
	assert.Equal(t, "a", spec.A)
	assert.Equal(t, "b", spec.B)
	assert.Equal(t, "c", spec.C)
}

func TestRouterRegisterDuplicateURL(t *testing.T) {
	var (
		container = nacelle.NewServiceContainer()
		logger    = nacelle.NewNilLogger()
		router    = NewRouter(container, logger)
	)

	err1 := router.Register("/users", &EmptySpec{})
	err2 := router.Register("/users", &EmptySpec{})

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "resource already registered")
}

func TestRouterRegisterWithMiddleware(t *testing.T) {
	var (
		container = nacelle.NewServiceContainer()
		logger    = nacelle.NewNilLogger()
		router    = NewRouter(container, logger)
		calls     = []string{}
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

	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/foo", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, calls)
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
