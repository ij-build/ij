package registry

import (
	"context"
	"io/ioutil"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type serverLogin struct {
	ctx      context.Context
	logger   logging.Logger
	env      environment.Environment
	registry *config.ServerRegistry
	runner   command.Runner
}

func NewServerLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.ServerRegistry,
) Login {
	return newServerLogin(
		ctx,
		logger,
		env,
		registry,
		command.NewRunner(logger),
	)
}

func newServerLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.ServerRegistry,
	runner command.Runner,
) Login {
	return &serverLogin{
		ctx:      ctx,
		logger:   logger,
		env:      env,
		registry: registry,
		runner:   runner,
	}
}

func (l *serverLogin) GetServer() (string, error) {
	return l.env.ExpandString(l.registry.Server)
}

func (l *serverLogin) Login() error {
	server, err := l.GetServer()
	if err != nil {
		return err
	}

	username, err := l.env.ExpandString(l.registry.Username)
	if err != nil {
		return err
	}

	password, err := getServerPassword(l.env, l.registry)
	if err != nil {
		return err
	}

	return login(
		l.ctx,
		l.runner,
		server,
		username,
		password,
	)
}

func getServerPassword(
	env environment.Environment,
	registry *config.ServerRegistry,
) (string, error) {
	passwordFile, err := env.ExpandString(registry.PasswordFile)
	if err != nil {
		return "", err
	}

	if passwordFile != "" {
		content, err := ioutil.ReadFile(passwordFile)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}

	return env.ExpandString(registry.Password)
}
