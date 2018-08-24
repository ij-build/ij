package runtime

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kballard/go-shellquote"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
)

type TaskBuilder struct {
	state         *State
	task          *config.Task
	containerName string
	env           environment.Environment
	args          []string
	command       []string
}

const (
	DefaultWorkspacePath = "/workspace"
	ScriptPath           = "/tmp/ij/script"
)

func NewTaskBuilder(
	state *State,
	task *config.Task,
	containerName string,
	env environment.Environment,
) *TaskBuilder {
	args := []string{
		"docker",
		"run",
		"--rm",
	}

	return &TaskBuilder{
		state:         state,
		task:          task,
		containerName: containerName,
		env:           env,
		args:          args,
	}
}

func (b *TaskBuilder) Build() ([]string, error) {
	builders := []func() error{
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

	for _, builder := range builders {
		if err := builder(); err != nil {
			return nil, err
		}
	}

	image, err := b.env.ExpandString(b.task.Image)
	if err != nil {
		return nil, err
	}

	return append(append(b.args, image), b.command...), nil
}

//
// Builders

func (b *TaskBuilder) addCommandOptions() error {
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

	b.addFlagValue("--entrypoint", entrypoint)
	b.command = commandArgs
	return nil
}

func (b *TaskBuilder) addContainerName() error {
	containerName, err := b.env.ExpandString(b.containerName)
	if err != nil {
		return err
	}

	b.addFlagValue("--name", containerName)
	return nil
}

func (b *TaskBuilder) addDetachOptions() error {
	if b.task.Detach {
		b.addFlag("-d")
	}

	return nil
}

func (b *TaskBuilder) addEnvironmentOptions() error {
	for _, line := range b.env.Serialize() {
		expanded, err := b.env.ExpandString(line)
		if err != nil {
			return err
		}

		b.addFlagValue("-e", expanded)
	}

	return nil
}

func (b *TaskBuilder) addHealthCheckOptions() error {
	hc := b.task.Healthcheck
	if hc == nil {
		return nil
	}

	command, err := b.env.ExpandString(hc.Command)
	if err != nil {
		return err
	}

	b.addFlagValue("--health-cmd", command)
	b.addFlagValue("--health-interval", hc.Interval.String())
	b.addFlagValue("--health-start-period", hc.StartPeriod.String())
	b.addFlagValue("--health-timeout", hc.Timeout.String())

	if hc.Retries > 0 {
		b.addFlagValue("--health-retries", fmt.Sprintf("%d", hc.Retries))
	}

	return nil
}

func (b *TaskBuilder) addLimitOptions() error {
	b.addFlagValue("--cpu-shares", b.state.cpuShares)
	b.addFlagValue("--memory", b.state.memory)
	return nil
}

func (b *TaskBuilder) addNetworkOptions() error {
	hostname, err := b.env.ExpandString(b.task.Hostname)
	if err != nil {
		return err
	}

	b.addFlagValue("--network", b.state.runID)
	b.addFlagValue("--network-alias", hostname)
	return nil
}

func (b *TaskBuilder) addScriptOptions() error {
	if b.task.Script == "" {
		return nil
	}

	script, err := b.env.ExpandString(b.task.Script)
	if err != nil {
		return err
	}

	path, err := b.state.scratch.WriteScript(script)
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

	b.addFlagValue("-v", mount)
	b.addFlagValue("--entrypoint", shell)
	b.command = []string{ScriptPath}
	return nil
}

func (b *TaskBuilder) addUserOptions() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	username, err := b.env.ExpandString(b.task.User)
	if err != nil {
		return err
	}

	b.addFlagValue("--user", username)
	b.addFlagValue("-e", fmt.Sprintf("UID=%s", user.Uid))
	b.addFlagValue("-e", fmt.Sprintf("GID=%s", user.Gid))
	return nil
}

func (b *TaskBuilder) addSSHOptions() error {
	if !b.state.enableSSHAgent {
		return nil
	}

	authSock := os.Getenv("SSH_AUTH_SOCK")
	b.addFlagValue("-e", "SSH_AUTH_SOCK")
	b.addFlagValue("-v", authSock+":"+authSock)
	return nil
}

func (b *TaskBuilder) addWorkspaceOptions() error {
	workspace, err := b.env.ExpandString(b.task.Workspace)
	if err != nil {
		return err
	}

	workspace, err = b.env.ExpandString(b.state.config.Workspace)
	if err != nil {
		return err
	}

	if workspace == "" {
		workspace = DefaultWorkspacePath
	}

	mount := fmt.Sprintf(
		"%s:%s",
		b.state.scratch.Workspace(),
		workspace,
	)

	b.addFlagValue("-v", mount)
	b.addFlagValue("-w", workspace)
	return nil
}

//
// Helpers

func (b *TaskBuilder) addFlag(flag string) {
	b.args = append(b.args, flag)
}

func (b *TaskBuilder) addFlagValue(flag, value string) {
	if value != "" {
		b.args = append(b.args, flag, value)
	}
}
