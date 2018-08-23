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
	state   *State
	task    *config.Task
	env     environment.Environment
	args    []string
	command []string
}

const (
	DefaultWorkspacePath = "/workspace"
	ScriptPath           = "/tmp/ij/script"
)

func NewTaskBuilder(
	state *State,
	containerName string,
	task *config.Task,
	env environment.Environment,
) *TaskBuilder {
	args := []string{
		"docker",
		"run",
		"--rm",
		"--name",
		containerName,
	}

	return &TaskBuilder{
		state: state,
		task:  task,
		env:   env,
		args:  args,
	}
}

// TODO - map environment

func (b *TaskBuilder) Build() ([]string, error) {
	builders := []func() error{
		b.addCommandOptions,
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

	return append(append(b.args, b.task.Image), b.command...), nil
}

//
// Builders

func (b *TaskBuilder) addCommandOptions() error {
	if b.task.Script != "" {
		return nil
	}

	command, err := shellquote.Split(b.task.Command)
	if err != nil {
		return err
	}

	b.addFlagValue("--entrypoint", b.task.Entrypoint)
	b.command = command
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
		b.addFlagValue("-e", line)
	}

	return nil
}

func (b *TaskBuilder) addHealthCheckOptions() error {
	hc := b.task.Healthcheck
	if hc == nil {
		return nil
	}

	b.addFlagValue("--health-cmd", hc.Command)
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
	b.addFlagValue("--network", b.state.runID)
	b.addFlagValue("--network-alias", b.task.Hostname)
	return nil
}

func (b *TaskBuilder) addScriptOptions() error {
	if b.task.Script == "" {
		return nil
	}

	path, err := b.state.scratch.WriteScript(b.task.Script)
	if err != nil {
		return err
	}

	mount := fmt.Sprintf(
		"%s:%s",
		path,
		ScriptPath,
	)

	shell := b.task.Shell
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

	b.addFlagValue("--user", b.task.User)
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
	workspacePath := b.task.WorkspacePath
	if workspacePath == "" {
		workspacePath = DefaultWorkspacePath
	}

	mount := fmt.Sprintf(
		"%s:%s",
		b.state.scratch.Workspace(),
		workspacePath,
	)

	b.addFlagValue("-v", mount)
	b.addFlagValue("-w", workspacePath)
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
