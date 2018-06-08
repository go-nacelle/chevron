package chevron

import (
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"
)

func Boot(name string, initializer RouteInitializer) {
	boostrapper := nacelle.NewBootstrapper(
		name,
		defaultSetupConfigs,
		setupProcessesFactory(initializer),
	)

	boostrapper.Boot()
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
