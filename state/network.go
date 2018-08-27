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
}

func NewNetwork(
	ctx context.Context,
	runID string,
	logger logging.Logger,
) (*Network, error) {
	n := &Network{
		ctx:    ctx,
		runID:  runID,
		logger: logger,
	}

	logger.Info(
		nil,
		"Creating network",
	)

	_, _, err := command.NewRunner(logger).RunForOutput(
		ctx,
		[]string{
			"docker",
			"network",
			"create",
			n.runID,
		},
	)

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Network) Teardown() {
	n.logger.Info(
		nil,
		"Removing network",
	)

	_, _, err := command.NewRunner(n.logger).RunForOutput(
		context.Background(),
		[]string{
			"docker",
			"network",
			"rm",
			n.runID,
		},
	)

	if err != nil {
		n.logger.Error(
			nil,
			"Failed to remove network: %s",
			err.Error(),
		)
	}
}
