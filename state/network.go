package state

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/logging"
)

type Network struct {
	ctx    context.Context
	runID  string
	logger logging.Logger
	runner command.Runner
}

func NewNetwork(
	ctx context.Context,
	runID string,
	logger logging.Logger,
) (*Network, error) {
	return newNetwork(
		ctx,
		runID,
		logger,
		command.NewRunner(logger),
	)
}

func newNetwork(
	ctx context.Context,
	runID string,
	logger logging.Logger,
	runner command.Runner,
) (*Network, error) {
	logger.Info(
		nil,
		"Creating network",
	)

	_, _, err := runner.RunForOutput(
		ctx,
		[]string{"docker", "network", "create", runID},
	)

	if err != nil {
		return nil, err
	}

	return &Network{
		ctx:    ctx,
		runID:  runID,
		logger: logger,
		runner: runner,
	}, nil
}

func (n *Network) Teardown() {
	n.logger.Info(
		nil,
		"Removing network",
	)

	_, _, err := n.runner.RunForOutput(
		context.Background(),
		[]string{"docker", "network", "rm", n.runID},
	)

	if err != nil {
		n.logger.Error(
			nil,
			"Failed to remove network: %s",
			err.Error(),
		)
	}
}
