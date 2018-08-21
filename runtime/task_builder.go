package runtime

import (
	"fmt"

	"github.com/kballard/go-shellquote"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
)

type TaskBuilder struct {
	runID         string
	containerName string
	workspace     *Workspace
	buildDir      *BuildDir
	task          *config.Task
	env           environment.Environment
	args          []string
}

const (
	// TODO - configure
	MountPoint  = "/workspace"
	ScriptMount = "/workspace/script"
)

func NewTaskBuilder(
	runID string,
	containerName string,
	workspace *Workspace,
	buildDir *BuildDir,
	task *config.Task,
	env environment.Environment,
) *TaskBuilder {
	return &TaskBuilder{
		runID:         runID,
		containerName: containerName,
		workspace:     workspace,
		buildDir:      buildDir,
		task:          task,
		env:           env,
		args:          []string{"docker", "run", "--rm"},
	}
}

//
// TODO - need to map environment as well
//

func (b *TaskBuilder) Build() ([]string, error) {
	b.addArgs("--name", b.containerName)
	b.addArgs("-w", MountPoint)
	b.addArgs("--network", b.runID)
	b.addArgs("-v", fmt.Sprintf(
		"%s:%s",
		b.workspace.VolumePath,
		MountPoint,
	))

	if b.task.Hostname != "" {
		b.addArgs("--network-alias", b.task.Hostname)
	}

	for _, line := range b.env.Serialize() {
		b.addArgs("-e", line)
	}

	// TODO - cpu shares
	// TODO - memory
	// TODO - UID/GID/user

	command, entrypoint, err := b.buildCommand()
	if err != nil {
		return nil, err
	}

	if entrypoint != "" {
		b.addArgs("--entrypoint", entrypoint)
	}

	if b.task.Detach {
		b.addArgs("-d")
	}

	b.buildHealthCheck()
	b.addArgs(b.task.Image)
	b.addArgs(command...)

	return b.args, nil
}

func (b *TaskBuilder) buildCommand() ([]string, string, error) {
	if b.task.Script == "" {
		command, err := shellquote.Split(b.task.Command)
		if err != nil {
			return nil, "", err
		}

		return command, b.task.Entrypoint, nil
	}

	path, err := b.buildDir.WriteScript(b.task.Script)
	if err != nil {
		return nil, "", err
	}

	b.addArgs("-v", fmt.Sprintf(
		"%s:%s",
		path,
		ScriptMount,
	))

	return []string{ScriptMount}, b.getShell(), nil
}

func (b *TaskBuilder) getShell() string {
	if b.task.Shell == "" {
		return "/bin/sh"
	}

	return b.task.Shell
}

func (b *TaskBuilder) buildHealthCheck() {
	healthcheck := b.task.Healthcheck
	if healthcheck == nil {
		return
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
}

func (b *TaskBuilder) addArgs(args ...string) {
	b.args = append(b.args, args...)
}
