package chevron

import (
	"net/http"

	"github.com/go-nacelle/httpbase"
	"github.com/go-nacelle/nacelle"
)

type (
	// ServerInitializer implements httpbase.ServerInitializer.
	ServerInitializer struct {
		Services    nacelle.ServiceContainer `service:"container"`
		Logger      nacelle.Logger           `service:"logger"`
		initializer RouteInitializer
		configs     []RouterConfigFunc
	}

	// RouteInitializer initializes a Router instance.
	RouteInitializer interface {
		// Init registers resources and middleware to the router.
		Init(config nacelle.Config, router Router) error
	}

	// RouteInitializerFunc is a function conforming to the RouteInitializer
	// interface.
	RouteInitializerFunc func(config nacelle.Config, router Router) error
)

// Init calls the wrapped function.
func (f RouteInitializerFunc) Init(config nacelle.Config, router Router) error {
	return f(config, router)
}

// NewInitializer creates a new ServerInitializer.
func NewInitializer(initializer RouteInitializer, configs ...RouterConfigFunc) httpbase.ServerInitializer {
	return &ServerInitializer{
		initializer: initializer,
		configs:     configs,
	}
}

// Init creates a router which becomes the server's handler and calls the
// attached route initializer.
func (i *ServerInitializer) Init(config nacelle.Config, server *http.Server) error {
	configs := append([]RouterConfigFunc{WithLogger(i.Logger)}, i.configs...)

	// TODO - control additional configs with env vars
	router := NewRouter(i.Services, i.Logger, configs...)
	server.Handler = router
	return i.initializer.Init(config, router)
}
