package config

import (
	"fmt"
	"time"
)

type (
	RunTask struct {
		TaskMeta
		Image                  string
		Command                string
		Shell                  string
		Script                 string
		Entrypoint             string
		User                   string
		Workspace              string
		Hostname               string
		Detach                 bool
		Healthcheck            *Healthcheck
		ExportEnvironmentFiles []string
	}

	Healthcheck struct {
		Command     string
		Interval    time.Duration
		Retries     int
		StartPeriod time.Duration
		Timeout     time.Duration
	}
)

func (t *RunTask) GetType() string {
	return "run"
}

func (t *RunTask) Extend(task Task) error {
	parent, ok := task.(*RunTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.extendMeta(parent.TaskMeta)
	t.Image = extendString(t.Image, parent.Image)
	t.Command = extendString(t.Command, parent.Command)
	t.Shell = extendString(t.Shell, parent.Shell)
	t.Script = extendString(t.Script, parent.Script)
	t.Entrypoint = extendString(t.Entrypoint, parent.Entrypoint)
	t.User = extendString(t.User, parent.User)
	t.Workspace = extendString(t.Workspace, parent.Workspace)
	t.Hostname = extendString(t.Hostname, parent.Hostname)
	t.Detach = extendBool(t.Detach, parent.Detach)
	t.Healthcheck.Extend(parent.Healthcheck)
	t.ExportEnvironmentFiles = append(parent.ExportEnvironmentFiles, t.ExportEnvironmentFiles...)
	return nil
}

func (t *Healthcheck) Extend(parent *Healthcheck) error {
	t.Command = extendString(t.Command, parent.Command)
	t.Interval = extendDuration(t.Interval, parent.Interval)
	t.Retries = extendInt(t.Retries, parent.Retries)
	t.StartPeriod = extendDuration(t.StartPeriod, parent.StartPeriod)
	t.Timeout = extendDuration(t.Timeout, parent.Timeout)
	return nil
}
