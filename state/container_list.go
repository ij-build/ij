package state

import (
	"context"
	"sync"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type ContainerList struct {
	description string
	target      func(string)
	logger      logging.Logger
	containers  map[string]struct{}
	mutex       sync.RWMutex
}

func NewContainerList(
	description string,
	target func(string),
	logger logging.Logger,
) *ContainerList {
	return &ContainerList{
		description: description,
		target:      target,
		logger:      logger,
		containers:  map[string]struct{}{},
	}
}

func (l *ContainerList) Add(containerName string) {
	l.mutex.Lock()
	l.containers[containerName] = struct{}{}
	l.mutex.Unlock()
}

func (l *ContainerList) Remove(containerName string) {
	l.mutex.Lock()
	delete(l.containers, containerName)
	l.mutex.Unlock()
}

func (l *ContainerList) Execute() {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	names := []string{}
	for containerName := range l.containers {
		names = append(names, containerName)
	}

	if len(names) == 0 {
		return
	}

	l.logger.Info(nil, l.description)
	util.RunParallelArgs(l.target, names...)
}

func NewContainerStopper(logger logging.Logger) *ContainerList {
	stopper := func(containerName string) {
		logger.Info(
			nil,
			"Stopping container %s",
			containerName,
		)

		args := []string{
			"docker",
			"kill",
			containerName,
		}

		// TODO - abstract this into state?
		_, _, err := command.NewRunner(logger).RunForOutput(
			context.Background(),
			args,
			nil,
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
		stopper,
		logger,
	)
}

func NewNetworkDisconnector(runID string, logger logging.Logger) *ContainerList {
	disconnect := func(containerName string) {
		logger.Info(
			nil,
			"Disconnecting container %s from network",
			containerName,
		)

		args := []string{
			"docker",
			"network",
			"disconnect",
			"--force",
			runID,
			containerName,
		}

		_, _, err := command.NewRunner(logger).RunForOutput(
			context.Background(),
			args,
			nil,
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
		disconnect,
		logger,
	)
}
