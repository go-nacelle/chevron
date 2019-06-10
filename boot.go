package chevron

import (
	"github.com/go-nacelle/httpbase"
	"github.com/go-nacelle/nacelle"
)

// BootAndExit creates a nacelle Bootstrapper with the given name and
// registers an HTTP server with the given route initializer. This method
// does not return.
func BootAndExit(name string, initializer RouteInitializer) {
	boostrapper := nacelle.NewBootstrapper(
		name,
		setupFactory(initializer),
	)

	boostrapper.BootAndExit()
}

func setupFactory(initializer RouteInitializer) func(nacelle.ProcessContainer, nacelle.ServiceContainer) error {
	return func(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
		processes.RegisterProcess(httpbase.NewServer(NewInitializer(initializer)))
		return nil
	}
}
