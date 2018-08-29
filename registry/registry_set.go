package registry

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type (
	RegistrySet struct {
		ctx         context.Context
		logger      logging.Logger
		env         environment.Environment
		runner      command.Runner
		namedLogins []*namedLogin
	}

	namedLogin struct {
		name  string
		login Login
	}

	loginFactory func(
		context.Context,
		logging.Logger,
		environment.Environment,
		config.Registry,
	) Login
)

func NewRegistrySet(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registries []config.Registry,
) (*RegistrySet, error) {
	return newRegistrySet(
		ctx,
		logger,
		env,
		registries,
		defaultLoginFactory,
		command.NewRunner(logger),
	)
}

func newRegistrySet(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registries []config.Registry,
	factory loginFactory,
	runner command.Runner,
) (*RegistrySet, error) {
	set := &RegistrySet{
		ctx:         ctx,
		logger:      logger,
		env:         env,
		runner:      runner,
		namedLogins: []*namedLogin{},
	}

	err := set.populateLoginMap(
		registries,
		factory,
	)

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (s *RegistrySet) Login() error {
	servers := []string{}
	for _, namedLogin := range s.namedLogins {
		s.logger.Info(
			nil,
			"Logging in to %s",
			namedLogin.name,
		)

		if err := namedLogin.login.Login(); err != nil {
			logoutRegistries(s.logger, s.runner, servers)
			return err
		}

		servers = append(servers, namedLogin.name)
	}

	return nil
}

func (s *RegistrySet) Logout() {
	if len(s.namedLogins) == 0 {
		return
	}

	s.logger.Info(
		nil,
		"Logging out of registries",
	)

	servers := []string{}
	for _, namedLogin := range s.namedLogins {
		servers = append(servers, namedLogin.name)
	}

	logoutRegistries(s.logger, s.runner, servers)
}

func (s *RegistrySet) populateLoginMap(
	registries []config.Registry,
	factory loginFactory,
) error {
	for _, registry := range registries {
		login := factory(s.ctx, s.logger, s.env, registry)

		server, err := login.GetServer()
		if err != nil {
			return err
		}

		s.namedLogins = append(s.namedLogins, &namedLogin{
			name:  server,
			login: login,
		})
	}

	return nil
}

//
// Helpers

func logoutRegistries(
	logger logging.Logger,
	runner command.Runner,
	servers []string,
) {
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

func defaultLoginFactory(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	reg config.Registry,
) Login {
	switch r := reg.(type) {
	case *config.ECRRegistry:
		return NewECRLogin(ctx, logger, env, r)
	case *config.GCRRegistry:
		return NewGCRLogin(ctx, logger, env, r)
	case *config.ServerRegistry:
		return NewServerLogin(ctx, logger, env, r)
	}

	panic("unexpected registry type")
}
