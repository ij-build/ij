package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/efritz/pvc/command"
)

type Workspace struct {
	runtime       *Runtime
	ContainerName string
	VolumePath    string
}

func NewWorkspace(runtime *Runtime) *Workspace {
	return &Workspace{
		runtime: runtime,
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
	containerName, err := command.RunForOutput(
		w.runtime.ctx,
		[]string{
			"docker",
			"create",
			fmt.Sprintf("--name=%s", w.runtime.id),
			"convey/workspace",
		},
	)

	if err != nil {
		return err
	}

	w.ContainerName = strings.TrimSpace(containerName)
	return nil
}

func (w *Workspace) inspect() error {
	volumePath, err := command.RunForOutput(
		w.runtime.ctx,
		[]string{
			"docker",
			"inspect",
			"--format",
			`{{range .Mounts}}{{if eq .Destination "/workspace"}}{{.Source}}{{end}}{{end}}`,
			w.ContainerName,
		},
	)

	if err != nil {
		w.Teardown()
		return err
	}

	w.VolumePath = strings.TrimSpace(volumePath)
	return nil
}
