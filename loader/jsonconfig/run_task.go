package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type (
	RunTask struct {
		Extends                string          `json:"extends"`
		Environment            json.RawMessage `json:"environment"`
		RequiredEnvironment    []string        `json:"required-environment"`
		Image                  string          `json:"image"`
		Command                string          `json:"command"`
		Shell                  string          `json:"shell"`
		Script                 string          `json:"script"`
		Entrypoint             string          `json:"entrypoint"`
		User                   string          `json:"user"`
		Workspace              string          `json:"workspace"`
		Hostname               string          `json:"hostname"`
		Detach                 bool            `json:"detach"`
		Healthcheck            *Healthcheck    `json:"healthcheck"`
		ExportEnvironmentFiles json.RawMessage `json:"export_environment_file"`
	}

	Healthcheck struct {
		Command     string   `json:"command"`
		Interval    Duration `json:"interval"`
		Retries     int      `json:"retries"`
		StartPeriod Duration `json:"start_period"`
		Timeout     Duration `json:"timeout"`
	}
)

func (t *RunTask) Translate(name string) (config.Task, error) {
	healthcheck, err := t.Healthcheck.Translate()
	if err != nil {
		return nil, err
	}

	environment, err := unmarshalStringList(t.Environment)
	if err != nil {
		return nil, err
	}

	exportedEnvironmentFiles, err := unmarshalStringList(t.ExportEnvironmentFiles)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:                name,
		Extends:             t.Extends,
		Environment:         environment,
		RequiredEnvironment: t.RequiredEnvironment,
	}

	return &config.RunTask{
		TaskMeta:               meta,
		Image:                  t.Image,
		Command:                t.Command,
		Shell:                  t.Shell,
		Script:                 t.Script,
		Entrypoint:             t.Entrypoint,
		User:                   t.User,
		Workspace:              t.Workspace,
		Hostname:               t.Hostname,
		Detach:                 t.Detach,
		Healthcheck:            healthcheck,
		ExportEnvironmentFiles: exportedEnvironmentFiles,
	}, nil
}

func (h *Healthcheck) Translate() (*config.Healthcheck, error) {
	if h == nil {
		return &config.Healthcheck{}, nil
	}

	return &config.Healthcheck{
		Command:     h.Command,
		Interval:    h.Interval.Duration,
		Retries:     h.Retries,
		StartPeriod: h.StartPeriod.Duration,
		Timeout:     h.Timeout.Duration,
	}, nil
}
