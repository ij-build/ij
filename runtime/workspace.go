package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/pvc/command"
	"github.com/efritz/pvc/logging"
)

type Workspace struct {
	runID         string
	ctx           context.Context
	logger        logging.Logger
	ContainerName string
	VolumePath    string
}

func NewWorkspace(runID string, ctx context.Context, logger logging.Logger) *Workspace {
	return &Workspace{
		runID:  runID,
		ctx:    ctx,
		logger: logger,
	}
}

func (w *Workspace) Setup() error {
	if err := w.create(); err != nil {
		return err
	}

	if err := w.inspect(); err != nil {
		w.Teardown()
		return err
	}

	return nil
}

func (w *Workspace) Teardown() error {
	w.logger.Info("Removing workspace")

	_, err := command.RunForOutput(
		context.Background(),
		[]string{
			"docker",
			"rm",
			"-v",
			w.ContainerName,
		},
	)

	return err
}

func (w *Workspace) create() error {
	args := []string{
		"docker",
		"create",
		fmt.Sprintf("--name=%s", w.runID),
		"convey/workspace",
	}

	w.logger.Info("Creating workspace")
	w.logger.Debug("Running command: %s", strings.Join(args, " "))

	containerName, err := command.RunForOutput(w.ctx, args)
	if err != nil {
		return err
	}

	w.ContainerName = strings.TrimSpace(containerName)
	return nil
}

func (w *Workspace) inspect() error {
	args := []string{
		"docker",
		"inspect",
		"--format",
		`{{range .Mounts}}{{if eq .Destination "/workspace"}}{{.Source}}{{end}}{{end}}`,
		w.ContainerName,
	}

	w.logger.Info("Inspecting workspace")
	w.logger.Debug("Running command: %s", strings.Join(args, " "))

	volumePath, err := command.RunForOutput(w.ctx, args)
	if err != nil {
		w.Teardown()
		return err
	}

	w.VolumePath = strings.TrimSpace(volumePath)
	return nil
}
