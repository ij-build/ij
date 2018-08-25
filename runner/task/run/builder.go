package run

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kballard/go-shellquote"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/state"
)

type builder struct {
	state         *state.State
	task          *config.RunTask
	containerName string
	env           environment.Environment
	command       []string
}

const (
	DefaultWorkspacePath = "/workspace"
	ScriptPath           = "/tmp/ij/script"
)

func build(
	state *state.State,
	task *config.RunTask,
	containerName string,
	env environment.Environment,
) ([]string, error) {
	b := &builder{
		state:         state,
		task:          task,
		containerName: containerName,
		env:           env,
	}

	builderFuncs := []command.BuildFunc{
		b.addCommandOptions,
		b.addContainerName,
		b.addDetachOptions,
		b.addEnvironmentOptions,
		b.addHealthCheckOptions,
		b.addLimitOptions,
		b.addNetworkOptions,
		b.addScriptOptions,
		b.addSSHOptions,
		b.addUserOptions,
		b.addWorkspaceOptions,
	}

	args := []string{
		"docker",
		"run",
		"--rm",
	}

	args, err := command.NewBuilder(builderFuncs, args).Build()
	if err != nil {
		return nil, err
	}

	image, err := b.env.ExpandString(b.task.Image)
	if err != nil {
		return nil, err
	}

	return append(append(args, image), b.command...), nil
}

//
// Builders

func (b *builder) addCommandOptions(cb *command.Builder) error {
	if b.task.Script != "" {
		return nil
	}

	command, err := b.env.ExpandString(b.task.Command)
	if err != nil {
		return err
	}

	entrypoint, err := b.env.ExpandString(b.task.Entrypoint)
	if err != nil {
		return err
	}

	commandArgs, err := shellquote.Split(command)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--entrypoint", entrypoint)
	b.command = commandArgs
	return nil
}

func (b *builder) addContainerName(cb *command.Builder) error {
	containerName, err := b.env.ExpandString(b.containerName)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--name", containerName)
	return nil
}

func (b *builder) addDetachOptions(cb *command.Builder) error {
	if b.task.Detach {
		cb.AddFlag("-d")
	}

	return nil
}

func (b *builder) addEnvironmentOptions(cb *command.Builder) error {
	for _, line := range b.env.Serialize() {
		expanded, err := b.env.ExpandString(line)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-e", expanded)
	}

	return nil
}

func (b *builder) addHealthCheckOptions(cb *command.Builder) error {
	command, err := b.env.ExpandString(b.task.Healthcheck.Command)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--health-cmd", command)
	cb.AddFlagValue("--health-interval", b.task.Healthcheck.Interval.String())
	cb.AddFlagValue("--health-start-period", b.task.Healthcheck.StartPeriod.String())
	cb.AddFlagValue("--health-timeout", b.task.Healthcheck.Timeout.String())

	if b.task.Healthcheck.Retries > 0 {
		cb.AddFlagValue("--health-retries", fmt.Sprintf(
			"%d",
			b.task.Healthcheck.Retries,
		))
	}

	return nil
}

func (b *builder) addLimitOptions(cb *command.Builder) error {
	cb.AddFlagValue("--cpu-shares", b.state.CPUShares)
	cb.AddFlagValue("--memory", b.state.Memory)
	return nil
}

func (b *builder) addNetworkOptions(cb *command.Builder) error {
	hostname, err := b.env.ExpandString(b.task.Hostname)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--network", b.state.RunID)
	cb.AddFlagValue("--network-alias", hostname)
	return nil
}

func (b *builder) addScriptOptions(cb *command.Builder) error {
	if b.task.Script == "" {
		return nil
	}

	script, err := b.env.ExpandString(b.task.Script)
	if err != nil {
		return err
	}

	path, err := b.state.Scratch.WriteScript(script)
	if err != nil {
		return err
	}

	mount := fmt.Sprintf(
		"%s:%s",
		path,
		ScriptPath,
	)

	shell, err := b.env.ExpandString(b.task.Shell)
	if err != nil {
		return err
	}

	if shell == "" {
		shell = "/bin/sh"
	}

	cb.AddFlagValue("-v", mount)
	cb.AddFlagValue("--entrypoint", shell)
	b.command = []string{ScriptPath}
	return nil
}

func (b *builder) addUserOptions(cb *command.Builder) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	username, err := b.env.ExpandString(b.task.User)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--user", username)
	cb.AddFlagValue("-e", fmt.Sprintf("UID=%s", user.Uid))
	cb.AddFlagValue("-e", fmt.Sprintf("GID=%s", user.Gid))
	return nil
}

func (b *builder) addSSHOptions(cb *command.Builder) error {
	if !b.state.EnableSSHAgent {
		return nil
	}

	authSock := os.Getenv("SSH_AUTH_SOCK")
	cb.AddFlagValue("-e", "SSH_AUTH_SOCK")
	cb.AddFlagValue("-v", authSock+":"+authSock)
	return nil
}

func (b *builder) addWorkspaceOptions(cb *command.Builder) error {
	workspace, err := b.env.ExpandString(b.task.Workspace)
	if err != nil {
		return err
	}

	workspace, err = b.env.ExpandString(b.state.Config.Workspace)
	if err != nil {
		return err
	}

	if workspace == "" {
		workspace = DefaultWorkspacePath
	}

	mount := fmt.Sprintf(
		"%s:%s",
		b.state.Scratch.Workspace(),
		workspace,
	)

	cb.AddFlagValue("-v", mount)
	cb.AddFlagValue("-w", workspace)
	return nil
}
