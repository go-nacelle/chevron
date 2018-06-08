package chevron

import (
	"net/http"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"
)

type (
	ServerInitializer struct {
		Container   *nacelle.ServiceContainer `service:"container"`
		Logger      nacelle.Logger            `service:"logger"`
		initializer RouteInitializer
	}

	RouteInitializer interface {
		Init(config nacelle.Config, router Router) error
	}

	RouteInitializerFunc func(config nacelle.Config, router Router) error
)

func (f RouteInitializerFunc) Init(config nacelle.Config, router Router) error {
	return f(config, router)
}

func NewInitializer(initializer RouteInitializer) process.HTTPServerInitializer {
	return &ServerInitializer{
		initializer: initializer,
	}
}

func (i *ServerInitializer) Init(config nacelle.Config, server *http.Server) error {
	// TODO - control additional configs with env vars
	router := NewRouter(i.Container, WithLogger(i.Logger))
	server.Handler = router
	return i.initializer.Init(config, router)
}
