package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kballard/go-shellquote"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/scratch"
	"github.com/efritz/ij/util"
)

type (
	RunTaskRunnerFactory func(
		*config.RunTask,
		environment.Environment,
		*logging.Prefix,
	) TaskRunner

	runTaskRunner struct {
		ctx              context.Context
		config           *config.Config
		runID            string
		scratch          *scratch.ScratchSpace
		containerLists   *ContainerLists
		containerOptions *containerOptions
		logger           logging.Logger
		loggerFactory    *logging.LoggerFactory
		task             *config.RunTask
		env              environment.Environment
		prefix           *logging.Prefix
	}

	containerOptions struct {
		EnableHostSSHAgent      bool
		EnableContainerSSHAgent bool
		CPUShares               string
		Memory                  string
	}

	runTaskCommandBuilderState struct {
		runID            string
		config           *config.Config
		containerOptions *containerOptions
		scratch          *scratch.ScratchSpace
		task             *config.RunTask
		containerName    string
		env              environment.Environment
	}
)

const (
	DefaultWorkspacePath = "/workspace"
	ScriptPath           = "/tmp/ij/script"
)

func NewRunTaskRunnerFactory(
	ctx context.Context,
	cfg *config.Config,
	runID string,
	scratch *scratch.ScratchSpace,
	containerLists *ContainerLists,
	containerOptions *containerOptions,
	logger logging.Logger,
	loggerFactory *logging.LoggerFactory,
) RunTaskRunnerFactory {
	return func(
		task *config.RunTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		return &runTaskRunner{
			ctx:              ctx,
			runID:            runID,
			config:           cfg,
			scratch:          scratch,
			containerLists:   containerLists,
			containerOptions: containerOptions,
			logger:           logger,
			loggerFactory:    loggerFactory,
			task:             task,
			env:              env,
			prefix:           prefix,
		}
	}
}

func (r *runTaskRunner) Run(context *RunContext) bool {
	r.logger.Info(
		r.prefix,
		"Beginning task",
	)

	containerName, err := util.MakeID()
	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to generate container id: %s",
			err.Error(),
		)

		return false
	}

	r.logger.Info(
		r.prefix,
		"Launching container %s",
		containerName,
	)

	builder, err := runTaskCommandBuilderFactory(
		r.runID,
		r.config,
		r.containerOptions,
		r.scratch,
		r.task,
		containerName,
		r.env,
	)

	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	args, _, err := builder.Build()
	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	if r.task.Detach {
		return r.runInBackground(containerName, args)
	}

	return r.runInForeground(context, containerName, args)
}

func (r *runTaskRunner) runInForeground(
	context *RunContext,
	containerName string,
	args []string,
) bool {
	logger, err := r.loggerFactory.Logger(
		r.prefix.Serialize(logging.NilColorPicker),
		false,
	)

	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to create task run log files: %s",
			err.Error(),
		)

		return false
	}

	r.containerLists.NetworkDisconnector.Add(containerName)
	defer r.containerLists.NetworkDisconnector.Remove(containerName)

	err = command.NewRunner(logger).Run(
		r.ctx,
		args,
		nil,
		r.prefix,
	)

	if err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	return r.exportEnvironmentFiles(context)
}

func (r *runTaskRunner) exportEnvironmentFiles(context *RunContext) bool {
	paths, err := r.env.ExpandSlice(r.task.ExportEnvironmentFiles)
	if err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Failed to build build export environment files: %s",
			err.Error(),
		)

		return false
	}

	for _, path := range paths {
		if !r.exportEnvironmentFile(context, path) {
			return false
		}
	}

	return true
}

func (r *runTaskRunner) exportEnvironmentFile(context *RunContext, path string) bool {
	workspace := r.scratch.Workspace()

	realPath, err := filepath.Abs(filepath.Join(workspace, path))
	if err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Failed to construct export environment file path: %s",
			err.Error(),
		)

		return false
	}

	if !strings.HasPrefix(realPath, workspace) {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"export environment file is outside of workspace directory: %s",
			realPath,
		)

		return false
	}

	r.logger.Info(
		r.prefix,
		"Injecting environment from file %s",
		fmt.Sprintf("~%s", realPath[len(workspace):]),
	)

	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to read environment file: %s",
			err.Error(),
		)

		return false
	}

	lines, err := environment.NormalizeEnvironmentFile(string(data))
	if err != nil {
		r.logger.Error(
			r.prefix,
			err.Error(),
		)

		return false
	}

	for _, line := range lines {
		context.ExportEnv(line)
	}

	return true
}

func (r *runTaskRunner) runInBackground(containerName string, args []string) bool {
	r.containerLists.ContainerStopper.Add(containerName)

	_, _, err := command.NewRunner(r.logger).RunForOutput(
		context.Background(),
		args,
		nil,
	)

	if err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	hasHealthcheck, err := hasHealthcheck(
		r.ctx,
		containerName,
		r.logger,
		r.prefix,
	)

	if err != nil {
		reportError(
			r.ctx,
			r.logger,
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

func (r *runTaskRunner) monitor(containerName string) bool {
	for {
		status, err := getHealthStatus(
			r.ctx,
			containerName,
			r.logger,
			r.prefix,
		)

		if err != nil {
			reportError(
				r.ctx,
				r.logger,
				r.prefix,
				"Failed to check container health: %s",
				err.Error(),
			)

			return false
		}

		if status == "healthy" {
			r.logger.Info(
				r.prefix,
				"Container is healthy",
			)

			return true
		}

		r.logger.Info(
			r.prefix,
			"Container is not yet healthy (currently %s)",
			status,
		)

		select {
		case <-time.After(r.config.Options.HealthcheckInterval):
		case <-r.ctx.Done():
			return false
		}
	}
}

func runTaskCommandBuilderFactory(
	runID string,
	config *config.Config,
	containerOptions *containerOptions,
	scratch *scratch.ScratchSpace,
	task *config.RunTask,
	containerName string,
	env environment.Environment,
) (*command.Builder, error) {
	s := &runTaskCommandBuilderState{
		runID:            runID,
		config:           config,
		containerOptions: containerOptions,
		scratch:          scratch,
		task:             task,
		containerName:    containerName,
		env:              env,
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

func (s *runTaskCommandBuilderState) addImageArg(cb *command.Builder) error {
	image, err := s.env.ExpandString(s.task.Image)
	if err != nil {
		return err
	}

	cb.AddArgs(image)
	return nil
}

func (s *runTaskCommandBuilderState) addCommandOptions(cb *command.Builder) error {
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

func (s *runTaskCommandBuilderState) addContainerName(cb *command.Builder) error {
	containerName, err := s.env.ExpandString(s.containerName)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--name", containerName)
	return nil
}

func (s *runTaskCommandBuilderState) addDetachOptions(cb *command.Builder) error {
	if s.task.Detach {
		cb.AddFlag("-d")
	}

	return nil
}

func (s *runTaskCommandBuilderState) addEnvironmentOptions(cb *command.Builder) error {
	for _, line := range s.env.Serialize() {
		expanded, err := s.env.ExpandString(line)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-e", expanded)
	}

	return nil
}

func (s *runTaskCommandBuilderState) addHealthcheckOptions(cb *command.Builder) error {
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

func (s *runTaskCommandBuilderState) addLimitOptions(cb *command.Builder) error {
	cb.AddFlagValue("--cpu-shares", s.containerOptions.CPUShares)
	cb.AddFlagValue("--memory", s.containerOptions.Memory)
	return nil
}

func (s *runTaskCommandBuilderState) addNetworkOptions(cb *command.Builder) error {
	hostname, err := s.env.ExpandString(s.task.Hostname)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--network", s.runID)
	cb.AddFlagValue("--network-alias", hostname)
	return nil
}

func (s *runTaskCommandBuilderState) addScriptOptions(cb *command.Builder) error {
	if s.task.Script == "" {
		return nil
	}

	script, err := s.env.ExpandString(s.task.Script)
	if err != nil {
		return err
	}

	path, err := s.scratch.WriteScript(script)
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

func (s *runTaskCommandBuilderState) addUserOptions(cb *command.Builder) error {
	uid, err := s.env.ExpandString("${UID}")
	if err != nil {
		return err
	}

	gid, err := s.env.ExpandString("${GID}")
	if err != nil {
		return err
	}

	username, err := s.env.ExpandString(s.task.User)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--user", username)
	cb.AddFlagValue("-e", fmt.Sprintf("UID=%s", uid))
	cb.AddFlagValue("-e", fmt.Sprintf("GID=%s", gid))
	return nil
}

func (s *runTaskCommandBuilderState) addSSHOptions(cb *command.Builder) error {
	if s.containerOptions.EnableHostSSHAgent {
		// Outside of the ssh-agent container we can just mount the host auth socket.
		authSock := os.Getenv("SSH_AUTH_SOCK")
		cb.AddFlagValue("-e", "SSH_AUTH_SOCK")
		cb.AddFlagValue("-v", fmt.Sprintf("%s:%s", authSock, authSock))
	}

	if s.containerOptions.EnableContainerSSHAgent {
		// Mount the socket from the ssh-agent container
		cb.AddFlagValue("-e", fmt.Sprintf("SSH_AUTH_SOCK=%s", SocketPath))
		cb.AddFlagValue("--volumes-from", fmt.Sprintf("%s-ssh-agent", s.runID))
	}

	return nil
}

func (s *runTaskCommandBuilderState) addWorkspaceOptions(cb *command.Builder) error {
	workspace, err := s.env.ExpandString(s.task.Workspace)
	if err != nil {
		return err
	}

	if workspace == "" {
		workspace, err = s.env.ExpandString(s.config.Workspace)
		if err != nil {
			return err
		}
	}

	if workspace == "" {
		workspace = DefaultWorkspacePath
	}

	mount := fmt.Sprintf(
		"%s:%s",
		s.scratch.Workspace(),
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
