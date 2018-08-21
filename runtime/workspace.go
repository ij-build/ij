package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/pvc/command"
	"github.com/efritz/pvc/logging"
)

type Workspace struct {
	ctx           context.Context
	runID         string
	logger        logging.Logger
	ContainerName string
	VolumePath    string
}

func NewWorkspace(
	ctx context.Context,
	runID string,
	logger logging.Logger,
) (*Workspace, error) {
	w := &Workspace{
		ctx:    ctx,
		runID:  runID,
		logger: logger,
	}

	if err := w.create(); err != nil {
		return nil, err
	}

	if err := w.inspect(); err != nil {
		w.Teardown()
		return nil, err
	}

	return w, nil
}

func (w *Workspace) Teardown() {
	w.logger.Info(
		nil,
		"Removing workspace",
	)

	_, _, err := command.RunForOutput(
		context.Background(),
		[]string{
			"docker",
			"rm",
			"-v",
			w.ContainerName,
		},
		w.logger,
	)

	if err != nil {
		w.logger.Error(
			nil,
			"Failed to remove workspace: %s",
			err.Error(),
		)
	}
}

func (w *Workspace) create() error {
	w.logger.Info(
		nil,
		"Creating workspace",
	)

	containerName, _, err := command.RunForOutput(
		w.ctx,
		[]string{
			"docker",
			"create",
			fmt.Sprintf("--name=%s", w.runID),
			"convey/workspace",
		},
		w.logger,
	)

	if err != nil {
		return err
	}

	w.ContainerName = strings.TrimSpace(containerName)
	return nil
}

func (w *Workspace) inspect() error {
	w.logger.Info(
		nil,
		"Inspecting workspace",
	)

	volumePath, _, err := command.RunForOutput(
		w.ctx,
		[]string{
			"docker",
			"inspect",
			"--format",
			`{{range .Mounts}}{{if eq .Destination "/workspace"}}{{.Source}}{{end}}{{end}}`,
			w.ContainerName,
		},
		w.logger,
	)

	if err != nil {
		w.Teardown()
		return err
	}

	w.VolumePath = strings.TrimSpace(volumePath)
	return nil
}
