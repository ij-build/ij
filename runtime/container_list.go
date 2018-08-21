package runtime

import (
	"context"
	"sync"

	"github.com/efritz/pvc/command"
	"github.com/efritz/pvc/logging"
	"github.com/efritz/pvc/util"
)

type ContainerList struct {
	logger      logging.Logger
	containers  []string
	mutex       sync.Mutex
	description string
	f           func(string)
}

func NewContainerList(
	description string,
	logger logging.Logger,
	f func(string),
) *ContainerList {
	return &ContainerList{
		description: description,
		f:           f,
		logger:      logger,
	}
}

func (l *ContainerList) RegisterContainer(containerName string) {
	l.mutex.Lock()
	l.containers = append(l.containers, containerName)
	l.mutex.Unlock()
}

func (l *ContainerList) Execute() {
	if len(l.containers) == 0 {
		return
	}

	l.logger.Info(nil, l.description)

	util.RunParallelArgs(l.f, l.containers...)
}

func NewContainerStopper(logger logging.Logger) *ContainerList {
	stopper := func(containerName string) {
		logger.Info(
			nil,
			"Stopping container %s",
			containerName,
		)

		_, _, err := command.RunForOutput(
			context.Background(),
			[]string{
				"docker",
				"kill",
				containerName,
			},
			logger,
		)

		if err != nil {
			logger.Error(
				nil,
				"Failed to stop container %s: %s",
				containerName,
				err.Error(),
			)

			return
		}

		logger.Info(
			nil,
			"Stopped container %s",
			containerName,
		)
	}

	return NewContainerList(
		"Stopping detached containers",
		logger,
		stopper,
	)
}

func NewNetworkDisconnector(runID string, logger logging.Logger) *ContainerList {
	disconnect := func(containerName string) {
		logger.Info(
			nil,
			"Disconnecting container %s from network",
			containerName,
		)

		_, _, err := command.RunForOutput(
			context.Background(),
			[]string{
				"docker",
				"network",
				"disconnect",
				"--force",
				runID,
				containerName,
			},
			logger,
		)

		if err != nil {
			logger.Error(
				nil,
				"Failed to disconnect container %s from network: %s",
				containerName,
				err.Error(),
			)

			return
		}

		logger.Info(
			nil,
			"Disconnected container %s from network",
			containerName,
		)
	}

	return NewContainerList(
		"Disconnecting containers from network",
		logger,
		disconnect,
	)
}
