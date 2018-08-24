package runtime

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
	"github.com/efritz/ij/util"
)

type TaskRunner struct {
	state  *State
	task   *config.Task
	prefix *logging.Prefix
	env    environment.Environment
}

func NewTaskRunner(
	state *State,
	task *config.Task,
	prefix *logging.Prefix,
	env environment.Environment,
) *TaskRunner {
	return &TaskRunner{
		state:  state,
		task:   task,
		prefix: prefix,
		env:    env,
	}
}

func (r *TaskRunner) Run() bool {
	r.state.logger.Info(
		r.prefix,
		"Beginning task",
	)

	ok, missing := util.ContainsAll(
		r.env.Keys(),
		r.task.RequiredEnvironment,
	)

	if !ok {
		r.state.logger.Error(
			r.prefix,
			"Missing environment values: %s",
			strings.Join(missing, ", "),
		)

		return false
	}

	containerName, err := util.MakeID()
	if err != nil {
		r.state.logger.Error(
			r.prefix,
			"Failed to generate container id: %s",
			err.Error(),
		)

		return false
	}

	r.state.logger.Info(
		r.prefix,
		"Launching container %s",
		containerName,
	)

	builder := NewTaskBuilder(
		r.state,
		r.task,
		containerName,
		r.env,
	)

	args, err := builder.Build()
	if err != nil {
		r.state.logger.Error(
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

func (r *TaskRunner) runInForeground(containerName string, args []string) bool {
	outfile, errfile, err := r.state.scratch.MakeLogFiles(
		r.prefix.Serialize(nil),
	)

	if err != nil {
		r.state.logger.Error(
			r.prefix,
			"Failed to create task run log files: %s",
			err.Error(),
		)

		return false
	}

	logger := r.state.logProcessor.Logger(
		outfile,
		errfile,
		false,
	)

	r.state.networkDisconnector.Add(containerName)
	defer r.state.networkDisconnector.Remove(containerName)

	commandErr := command.Run(
		r.state.ctx,
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

func (r *TaskRunner) exportEnvironmentFiles() bool {
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

func (r *TaskRunner) exportEnvironmentFile(path string) bool {
	realPath, err := filepath.Abs(filepath.Join(
		r.state.scratch.Workspace(),
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

	workspace := r.state.scratch.Workspace()

	if !strings.HasPrefix(realPath, workspace) {
		r.state.ReportError(
			r.prefix,
			"export environment file is outside of workspace directory: %s",
			realPath,
		)

		return false
	}

	r.state.logger.Info(
		r.prefix,
		"Injecting environment from file %s",
		fmt.Sprintf("~%s", realPath[len(workspace):]),
	)

	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		r.state.logger.Error(
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
			r.state.logger.Error(
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

func (r *TaskRunner) runInBackground(containerName string, args []string) bool {
	r.state.containerStopper.Add(containerName)

	_, _, err := command.RunForOutput(
		context.Background(),
		args,
		r.state.logger,
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
		r.state.ctx,
		containerName,
		r.state.logger,
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

func (r *TaskRunner) monitor(containerName string) bool {
	for {
		status, err := getHealthStatus(
			r.state.ctx,
			containerName,
			r.state.logger,
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
			r.state.logger.Info(
				r.prefix,
				"Container is healthy",
			)

			return true
		}

		r.state.logger.Info(
			r.prefix,
			"Container is not yet healthy (currently %s)",
			status,
		)

		select {
		case <-time.After(r.state.healthcheckInterval):
		case <-r.state.ctx.Done():
			return false
		}
	}
}
