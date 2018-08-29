package state

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry"
)

type (
	RegistryList struct {
		logger  logging.Logger
		runner  command.Runner
		servers []string
	}

	loginFactory func(
		context.Context,
		logging.Logger,
		environment.Environment,
		config.Registry,
	) registry.Login
)

func NewRegistryList(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registries []config.Registry,
) (*RegistryList, error) {
	return newRegistryList(
		ctx,
		logger,
		env,
		registries,
		buildLogin,
		command.NewRunner(logger),
	)
}

func newRegistryList(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registries []config.Registry,
	factory loginFactory,
	runner command.Runner,
) (*RegistryList, error) {
	servers := []string{}
	for _, registry := range registries {
		logger.Info(
			nil,
			"Logging into %s registry",
			registry.GetType(),
		)

		server, err := factory(ctx, logger, env, registry).Login()
		if err != nil {
			logoutRegistries(logger, runner, servers)
			return nil, err
		}

		servers = append(servers, server)
	}

	return &RegistryList{
		logger:  logger,
		runner:  runner,
		servers: servers,
	}, nil
}

func (l *RegistryList) Teardown() {
	if len(l.servers) == 0 {
		return
	}

	l.logger.Info(
		nil,
		"Logging out of registries",
	)

	logoutRegistries(l.logger, l.runner, l.servers)
}

//
// Helpers

func buildLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	reg config.Registry,
) registry.Login {
	switch r := reg.(type) {
	case *config.ECRRegistry:
		return registry.NewECRLogin(ctx, logger, env, r)
	case *config.GCRRegistry:
		return registry.NewGCRLogin(ctx, logger, env, r)
	case *config.ServerRegistry:
		return registry.NewServerLogin(ctx, logger, env, r)
	}

	panic("unexpected registry type")
}

func logoutRegistries(logger logging.Logger, runner command.Runner, servers []string) {
	for _, server := range servers {
		logger.Info(
			nil,
			"Logging out of %s",
			server,
		)

		args := []string{
			"docker",
			"logout",
			server,
		}

		err := runner.Run(
			context.Background(),
			args,
			nil,
			nil,
		)

		if err != nil {
			logger.Error(
				nil,
				"Failed to log out of registry: %s",
				err.Error(),
			)
		}
	}
}
