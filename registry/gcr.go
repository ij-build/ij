package registry

import (
	"context"
	"io/ioutil"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type gcrLogin struct {
	ctx      context.Context
	logger   logging.Logger
	env      environment.Environment
	registry *config.GCRRegistry
	runner   command.Runner
}

func NewGCRLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.GCRRegistry,
) Login {
	return newGCRLogin(
		ctx,
		logger,
		env,
		registry,
		command.NewRunner(logger),
	)
}

func newGCRLogin(
	ctx context.Context,
	logger logging.Logger,
	env environment.Environment,
	registry *config.GCRRegistry,
	runner command.Runner,
) Login {
	return &gcrLogin{
		ctx:      ctx,
		logger:   logger,
		env:      env,
		registry: registry,
		runner:   runner,
	}
}

func (l *gcrLogin) Login() (string, error) {
	server := "https://gcr.io"

	password, err := getGCRPassword(l.env, l.registry)
	if err != nil {
		return "", err
	}

	err = login(
		l.ctx,
		l.runner,
		server,
		"_json_key",
		string(password),
	)

	if err != nil {
		return "", err
	}

	return server, nil
}

func getGCRPassword(
	env environment.Environment,
	registry *config.GCRRegistry,
) (string, error) {
	keyFile, err := env.ExpandString(registry.KeyFile)
	if err != nil {
		return "", err
	}

	password, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", err
	}

	return string(password), nil
}
