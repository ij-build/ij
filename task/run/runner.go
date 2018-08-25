package run

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
	"github.com/efritz/ij/util"
)

type Runner struct {
	state  *state.State
	task   *config.Task
	prefix *logging.Prefix
	env    environment.Environment
}

func NewRunner(
	state *state.State,
	task *config.Task,
	prefix *logging.Prefix,
	env environment.Environment,
) *Runner {
	return &Runner{
		state:  state,
		task:   task,
		prefix: prefix,
		env:    env,
	}
}

func (r *Runner) Run() bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning task",
	)

	ok, missing := util.ContainsAll(
		r.env.Keys(),
		r.task.RequiredEnvironment,
	)

	if !ok {
		r.state.Logger.Error(
			r.prefix,
			"Missing environment values: %s",
			strings.Join(missing, ", "),
		)

		return false
	}

	containerName, err := util.MakeID()
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to generate container id: %s",
			err.Error(),
		)

		return false
	}

	r.state.Logger.Info(
		r.prefix,
		"Launching container %s",
		containerName,
	)

	builder := NewBuilder(
		r.state,
		r.task,
		containerName,
		r.env,
	)

	args, err := builder.Build()
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	if !r.task.Detach {
		return r.runInForeground(containerName, args)
	}

	return r.runInBackground(containerName, args)
}

func (r *Runner) runInForeground(containerName string, args []string) bool {
	outfile, errfile, err := r.state.Scratch.MakeLogFiles(
		r.prefix.Serialize(nil),
	)

	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to create task run log files: %s",
			err.Error(),
		)

		return false
	}

	logger := r.state.LogProcessor.Logger(
		outfile,
		errfile,
		false,
	)

	r.state.NetworkDisconnector.Add(containerName)
	defer r.state.NetworkDisconnector.Remove(containerName)

	commandErr := command.Run(
		r.state.Context,
		args,
		logger,
		r.prefix,
	)

	if commandErr != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			commandErr.Error(),
		)

		return false
	}

	return r.exportEnvironmentFiles()
}

func (r *Runner) exportEnvironmentFiles() bool {
	paths, err := r.env.ExpandSlice(r.task.ExportEnvironmentFiles)
	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Failed to build build export environment files: %s",
			err.Error(),
		)

		return false
	}

	for _, path := range paths {
		if !r.exportEnvironmentFile(path) {
			return false
		}
	}

	return true
}

func (r *Runner) exportEnvironmentFile(path string) bool {
	realPath, err := filepath.Abs(filepath.Join(
		r.state.Scratch.Workspace(),
		path,
	))

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Failed to construct export environment file path: %s",
			err.Error(),
		)

		return false
	}

	workspace := r.state.Scratch.Workspace()

	if !strings.HasPrefix(realPath, workspace) {
		r.state.ReportError(
			r.prefix,
			"export environment file is outside of workspace directory: %s",
			realPath,
		)

		return false
	}

	r.state.Logger.Info(
		r.prefix,
		"Injecting environment from file %s",
		fmt.Sprintf("~%s", realPath[len(workspace):]),
	)

	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to read environment file: %s",
			err.Error(),
		)

		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || line[0] == '#' {
			continue
		}

		if !strings.Contains(line, "=") {
			r.state.Logger.Error(
				r.prefix,
				"Malformed entry in environments file: %s",
				line,
			)

			return false
		}

		r.state.ExportEnv(line)
	}

	return true
}

func (r *Runner) runInBackground(containerName string, args []string) bool {
	r.state.ContainerStopper.Add(containerName)

	_, _, err := command.RunForOutput(
		context.Background(),
		args,
		r.state.Logger,
	)

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	hasHealthcheck, err := hasHealthCheck(
		r.state.Context,
		containerName,
		r.state.Logger,
		r.prefix,
	)

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Failed to determine if container has a healthcheck: %s",
			err.Error(),
		)

		return false
	}

	if !hasHealthcheck {
		return true
	}

	return r.monitor(containerName)
}

func (r *Runner) monitor(containerName string) bool {
	for {
		status, err := getHealthStatus(
			r.state.Context,
			containerName,
			r.state.Logger,
			r.prefix,
		)

		if err != nil {
			r.state.ReportError(
				r.prefix,
				"Failed to check container health: %s",
				err.Error(),
			)

			return false
		}

		if status == "healthy" {
			r.state.Logger.Info(
				r.prefix,
				"Container is healthy",
			)

			return true
		}

		r.state.Logger.Info(
			r.prefix,
			"Container is not yet healthy (currently %s)",
			status,
		)

		select {
		case <-time.After(r.state.HealthcheckInterval):
		case <-r.state.Context.Done():
			return false
		}
	}
}
