package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/kballard/go-shellquote"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
	"github.com/efritz/ij/util"
)

type (
	runCommandRunner struct {
		state  *state.State
		task   *config.RunTask
		prefix *logging.Prefix
		env    environment.Environment
	}

	runCommandBuilderState struct {
		state         *state.State
		task          *config.RunTask
		containerName string
		env           environment.Environment
	}
)

const (
	DefaultWorkspacePath = "/workspace"
	ScriptPath           = "/tmp/ij/script"
)

func NewRunCommandRunner(
	state *state.State,
	runTask *config.RunTask,
	prefix *logging.Prefix,
	env environment.Environment,
) Runner {
	return &runCommandRunner{
		state:  state,
		task:   runTask,
		prefix: prefix,
		env:    env,
	}
}

func (r *runCommandRunner) Run() bool {
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

	builder, err := runCommandBuilderFactory(
		r.state,
		r.task,
		containerName,
		r.env,
	)

	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	args, _, err := builder.Build()
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

func (r *runCommandRunner) runInForeground(containerName string, args []string) bool {
	outfile, errfile, err := r.state.Scratch.MakeLogFiles(
		r.prefix.Serialize(logging.NilColorPicker),
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

	err = command.NewRunner(logger).Run(
		r.state.Context,
		args,
		nil,
		r.prefix,
	)

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	return r.exportEnvironmentFiles()
}

func (r *runCommandRunner) exportEnvironmentFiles() bool {
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

func (r *runCommandRunner) exportEnvironmentFile(path string) bool {
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

func (r *runCommandRunner) runInBackground(containerName string, args []string) bool {
	r.state.ContainerStopper.Add(containerName)

	_, _, err := command.NewRunner(r.state.Logger).RunForOutput(
		context.Background(),
		args,
		nil,
	)

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	hasHealthcheck, err := hasHealthcheck(
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

func (r *runCommandRunner) monitor(containerName string) bool {
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
		case <-time.After(r.state.Config.HealthcheckInterval):
		case <-r.state.Context.Done():
			return false
		}
	}
}

func runCommandBuilderFactory(
	state *state.State,
	task *config.RunTask,
	containerName string,
	env environment.Environment,
) (*command.Builder, error) {
	s := &runCommandBuilderState{
		state:         state,
		task:          task,
		containerName: containerName,
		env:           env,
	}

	return command.NewBuilder(
		[]string{
			"docker",
			"run",
			"--rm",
		},
		[]command.BuildFunc{
			s.addImageArg,
			s.addCommandOptions, // Populates command, must come second
			s.addScriptOptions,  // Populates command, must come second
			s.addContainerName,
			s.addDetachOptions,
			s.addEnvironmentOptions,
			s.addHealthcheckOptions,
			s.addLimitOptions,
			s.addNetworkOptions,
			s.addSSHOptions,
			s.addUserOptions,
			s.addWorkspaceOptions,
		},
	), nil
}

//
// Builders

func (s *runCommandBuilderState) addImageArg(cb *command.Builder) error {
	image, err := s.env.ExpandString(s.task.Image)
	if err != nil {
		return err
	}

	cb.AddArgs(image)
	return nil
}

func (s *runCommandBuilderState) addCommandOptions(cb *command.Builder) error {
	if s.task.Script != "" {
		return nil
	}

	command, err := s.env.ExpandString(s.task.Command)
	if err != nil {
		return err
	}

	entrypoint, err := s.env.ExpandString(s.task.Entrypoint)
	if err != nil {
		return err
	}

	commandArgs, err := shellquote.Split(command)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--entrypoint", entrypoint)
	cb.AddArgs(commandArgs...)
	return nil
}

func (s *runCommandBuilderState) addContainerName(cb *command.Builder) error {
	containerName, err := s.env.ExpandString(s.containerName)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--name", containerName)
	return nil
}

func (s *runCommandBuilderState) addDetachOptions(cb *command.Builder) error {
	if s.task.Detach {
		cb.AddFlag("-d")
	}

	return nil
}

func (s *runCommandBuilderState) addEnvironmentOptions(cb *command.Builder) error {
	for _, line := range s.env.Serialize() {
		expanded, err := s.env.ExpandString(line)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-e", expanded)
	}

	return nil
}

func (s *runCommandBuilderState) addHealthcheckOptions(cb *command.Builder) error {
	command, err := s.env.ExpandString(s.task.Healthcheck.Command)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--health-cmd", command)
	cb.AddFlagValue("--health-interval", s.task.Healthcheck.Interval.String())
	cb.AddFlagValue("--health-start-period", s.task.Healthcheck.StartPeriod.String())
	cb.AddFlagValue("--health-timeout", s.task.Healthcheck.Timeout.String())

	if s.task.Healthcheck.Retries > 0 {
		cb.AddFlagValue("--health-retries", fmt.Sprintf(
			"%d",
			s.task.Healthcheck.Retries,
		))
	}

	return nil
}

func (s *runCommandBuilderState) addLimitOptions(cb *command.Builder) error {
	cb.AddFlagValue("--cpu-shares", s.state.CPUShares)
	cb.AddFlagValue("--memory", s.state.Memory)
	return nil
}

func (s *runCommandBuilderState) addNetworkOptions(cb *command.Builder) error {
	hostname, err := s.env.ExpandString(s.task.Hostname)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--network", s.state.RunID)
	cb.AddFlagValue("--network-alias", hostname)
	return nil
}

func (s *runCommandBuilderState) addScriptOptions(cb *command.Builder) error {
	if s.task.Script == "" {
		return nil
	}

	script, err := s.env.ExpandString(s.task.Script)
	if err != nil {
		return err
	}

	path, err := s.state.Scratch.WriteScript(script)
	if err != nil {
		return err
	}

	mount := fmt.Sprintf(
		"%s:%s",
		path,
		ScriptPath,
	)

	shell, err := s.env.ExpandString(s.task.Shell)
	if err != nil {
		return err
	}

	if shell == "" {
		shell = "/bin/sh"
	}

	cb.AddFlagValue("-v", mount)
	cb.AddFlagValue("--entrypoint", shell)
	cb.AddArgs(ScriptPath)
	return nil
}

func (s *runCommandBuilderState) addUserOptions(cb *command.Builder) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	username, err := s.env.ExpandString(s.task.User)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--user", username)
	cb.AddFlagValue("-e", fmt.Sprintf("UID=%s", user.Uid))
	cb.AddFlagValue("-e", fmt.Sprintf("GID=%s", user.Gid))
	return nil
}

func (s *runCommandBuilderState) addSSHOptions(cb *command.Builder) error {
	if !s.state.EnableSSHAgent {
		return nil
	}

	authSock := os.Getenv("SSH_AUTH_SOCK")
	cb.AddFlagValue("-e", "SSH_AUTH_SOCK")
	cb.AddFlagValue("-v", authSock+":"+authSock)
	return nil
}

func (s *runCommandBuilderState) addWorkspaceOptions(cb *command.Builder) error {
	workspace, err := s.env.ExpandString(s.task.Workspace)
	if err != nil {
		return err
	}

	workspace, err = s.env.ExpandString(s.state.Config.Workspace)
	if err != nil {
		return err
	}

	if workspace == "" {
		workspace = DefaultWorkspacePath
	}

	mount := fmt.Sprintf(
		"%s:%s",
		s.state.Scratch.Workspace(),
		workspace,
	)

	cb.AddFlagValue("-v", mount)
	cb.AddFlagValue("-w", workspace)
	return nil
}

//
// Helpers

func hasHealthcheck(
	ctx context.Context,
	containerName string,
	logger logging.Logger,
	prefix *logging.Prefix,
) (bool, error) {
	logger.Debug(prefix, "Checking if container has a healthcheck")

	args := []string{
		"docker",
		"inspect",
		"-f",
		"{{if .Config.Healthcheck}}true{{else}}false{{end}}",
		containerName,
	}

	out, _, err := command.NewRunner(logger).RunForOutput(
		ctx,
		args,
		nil,
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

	args := []string{
		"docker",
		"inspect",
		"-f",
		"{{.State.Health.Status}}",
		containerName,
	}

	out, _, err := command.NewRunner(logger).RunForOutput(
		ctx,
		args,
		nil,
	)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
