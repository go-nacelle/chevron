package chevron

import (
	"os"

	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"
)

// Boot creates a nacelle Bootstrapper with the given name and registers
// an HTTP server with the given route initializer. This method does not
// return.
func Boot(name string, initializer RouteInitializer) {
	boostrapper := nacelle.NewBootstrapper(
		name,
		defaultSetupConfigs,
		setupProcessesFactory(initializer),
	)

	// TODO - check for uses in other projects
	os.Exit(boostrapper.Boot())
}

func defaultSetupConfigs(config nacelle.Config) error {
	return config.Register(process.HTTPConfigToken, &process.HTTPConfig{})
}

func setupProcessesFactory(initializer RouteInitializer) func(*nacelle.ProcessRunner, *nacelle.ServiceContainer) error {
	return func(runner *nacelle.ProcessRunner, container *nacelle.ServiceContainer) error {
		runner.RegisterProcess(process.NewHTTPServer(NewInitializer(initializer)))
		return nil
	}
}
