package subcommand

import (
	"context"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry"
)

func NewLoginCommand(appOptions *AppOptions) CommandRunner {
	return func(config *config.Config) error {
		return withRegistrySet(config, appOptions, func(registrySet *registry.RegistrySet, logger logging.Logger) error {
			if err := registrySet.Login(); err != nil {
				return fmt.Errorf(
					"failed to log in to registries: %s",
					err.Error(),
				)
			}

			return nil
		})
	}
}

func NewLogoutCommand(appOptions *AppOptions) CommandRunner {
	return func(config *config.Config) error {
		return withRegistrySet(config, appOptions, func(registrySet *registry.RegistrySet, logger logging.Logger) error {
			registrySet.Logout()
			return nil
		})
	}
}

func withRegistrySet(
	config *config.Config,
	appOptions *AppOptions,
	f func(*registry.RegistrySet, logging.Logger) error,
) error {
	logProcessor := logging.NewProcessor(appOptions.Verbose, appOptions.Colorize)
	logProcessor.Start()
	defer logProcessor.Shutdown()

	logger := logProcessor.Logger(
		logging.NilWriter,
		logging.NilWriter,
		true,
	)

	registryEnv := environment.Merge(
		environment.New(config.Environment),
		environment.New(appOptions.Env),
	)

	registrySet, err := registry.NewRegistrySet(
		context.Background(),
		logger,
		registryEnv,
		config.Registries,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create registry set: %s",
			err.Error(),
		)
	}

	return f(registrySet, logger)
}
