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
	// TODO - configure
	MountPoint  = "/workspace"
	ScriptMount = "/workspace/script"
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

	if b.task.Entrypoint != "" {
		b.addArgs("--entrypoint", b.task.Entrypoint)
	}

	b.command = command
	return nil
}

func (b *TaskBuilder) addDetachOptions() error {
	if b.task.Detach {
		b.addArgs("-d")
	}

	return nil
}

func (b *TaskBuilder) addEnvironmentOptions() error {
	for _, line := range b.env.Serialize() {
		b.addArgs("-e", line)
	}

	return nil
}

func (b *TaskBuilder) addHealthCheckOptions() error {
	healthcheck := b.task.Healthcheck
	if healthcheck == nil {
		return nil
	}

	if healthcheck.Command != "" {
		b.addArgs("--health-cmd", healthcheck.Command)
	}

	if healthcheck.Interval.Duration > 0 {
		b.addArgs("--health-interval", healthcheck.Interval.String())
	}

	if healthcheck.Retries > 0 {
		b.addArgs("--health-retries", fmt.Sprintf("%d", healthcheck.Retries))
	}

	if healthcheck.StartPeriod.Duration > 0 {
		b.addArgs("--health-start-period", healthcheck.StartPeriod.String())
	}

	if healthcheck.Timeout.Duration > 0 {
		b.addArgs("--health-timeout", healthcheck.Timeout.String())
	}

	return nil
}

func (b *TaskBuilder) addLimitOptions() error {
	if b.state.cpuShares != "" {
		b.addArgs("--cpu-shares", b.state.cpuShares)
	}

	if b.state.memory != "" {
		b.addArgs("--memory", b.state.memory)
	}

	return nil
}

func (b *TaskBuilder) addNetworkOptions() error {
	b.addArgs("--network", b.state.runID)

	if b.task.Hostname != "" {
		b.addArgs("--network-alias", b.task.Hostname)
	}

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

	b.addArgs("-v", fmt.Sprintf(
		"%s:%s",
		path,
		ScriptMount,
	))

	if b.task.Shell == "" {
		b.addArgs("--entrypoint", "/bin/sh")
	} else {
		b.addArgs("--entrypoint", b.task.Shell)
	}

	b.command = []string{ScriptMount}
	return nil
}

func (b *TaskBuilder) addUserOptions() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	b.addArgs("-e", fmt.Sprintf("UID=%s", user.Uid))
	b.addArgs("-e", fmt.Sprintf("GID=%s", user.Gid))

	if b.task.User != "" {
		b.addArgs("--user", b.task.User)
	}

	return nil
}

func (b *TaskBuilder) addSSHOptions() error {
	if !b.state.enableSSHAgent {
		return nil
	}

	authSock := os.Getenv("SSH_AUTH_SOCK")
	b.addArgs("-e", "SSH_AUTH_SOCK")
	b.addArgs("-v", authSock+":"+authSock)
	return nil
}

func (b *TaskBuilder) addWorkspaceOptions() error {
	b.addArgs("-w", MountPoint)
	b.addArgs("-v", fmt.Sprintf(
		"%s:%s",
		b.state.scratch.Workspace(),
		MountPoint,
	))

	return nil
}

//
// Helpers

func (b *TaskBuilder) addArgs(args ...string) {
	b.args = append(b.args, args...)
}
