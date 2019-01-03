package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	RunTask struct {
		TaskMeta
		Image                  string       `json:"image,omitempty"`
		Command                string       `json:"command,omitempty"`
		Shell                  string       `json:"shell,omitempty"`
		Script                 string       `json:"script,omitempty"`
		Entrypoint             string       `json:"entrypoint,omitempty"`
		User                   string       `json:"user,omitempty"`
		Workspace              string       `json:"workspace,omitempty"`
		Hostname               string       `json:"hostname,omitempty"`
		Detach                 bool         `json:"detach,omitempty"`
		Healthcheck            *Healthcheck `json:"healthcheck,omitempty"`
		ExportEnvironmentFiles []string     `json:"export-environment-files,omitempty"`
	}

	// Note: Healthcheck must serialize itself manually due to the time.Duration fields.

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

func (h *Healthcheck) Extend(parent *Healthcheck) error {
	h.Command = extendString(h.Command, parent.Command)
	h.Interval = extendDuration(h.Interval, parent.Interval)
	h.Retries = extendInt(h.Retries, parent.Retries)
	h.StartPeriod = extendDuration(h.StartPeriod, parent.StartPeriod)
	h.Timeout = extendDuration(h.Timeout, parent.Timeout)
	return nil
}

func (t *RunTask) MarshalJSON() ([]byte, error) {
	type Alias RunTask

	return json.Marshal(&struct {
		*Alias
		Type string `json:"type"`
	}{
		Alias: (*Alias)(t),
		Type:  t.GetType(),
	})
}

func (h *Healthcheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Command     string `json:"command,omitempty"`
		Interval    string `json:"interval,omitempty"`
		Retries     int    `json:"retries,omitempty"`
		StartPeriod string `json:"start-period,omitempty"`
		Timeout     string `json:"timeout,omitempty"`
	}{
		Command:     h.Command,
		Interval:    durationString(h.Interval),
		Retries:     h.Retries,
		StartPeriod: durationString(h.StartPeriod),
		Timeout:     durationString(h.Timeout),
	})
}
