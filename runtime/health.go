package runtime

import (
	"context"
	"strings"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/logging"
)

func hasHealthCheck(
	ctx context.Context,
	containerName string,
	logger logging.Logger,
	prefix *logging.Prefix,
) (bool, error) {
	logger.Debug(prefix, "Checking if container has a healthcheck")

	out, _, err := command.RunForOutput(
		ctx,
		[]string{
			"docker",
			"inspect",
			"-f",
			"{{if .Config.Healthcheck}}true{{else}}false{{end}}",
			containerName,
		},
		logger,
	)

	if err != nil {
		return false, err
	}

	return strings.TrimSpace(out) == "true", nil
}

func getHealthStatus(
	ctx context.Context,
	containerName string,
	logger logging.Logger,
	prefix *logging.Prefix,
) (string, error) {
	logger.Debug(prefix, "Checking container health")

	out, _, err := command.RunForOutput(
		ctx,
		[]string{
			"docker",
			"inspect",
			"-f",
			"{{.State.Health.Status}}",
			containerName,
		},
		logger,
	)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
