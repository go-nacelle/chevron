package chevron

import (
	"os"

	"github.com/efritz/nacelle"
	basehttp "github.com/efritz/nacelle/base/http"
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
	return config.Register(basehttp.ConfigToken, &basehttp.Config{})
}

func setupProcessesFactory(initializer RouteInitializer) func(nacelle.ProcessContainer, nacelle.ServiceContainer) error {
	return func(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
		processes.RegisterProcess(basehttp.NewServer(NewInitializer(initializer)))
		return nil
	}
}
