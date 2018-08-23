package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	Task struct {
		Name                string
		Extends             string       `json:"extends"`
		Image               string       `json:"image"`
		Command             string       `json:"command"`
		Shell               string       `json:"shell"`
		Script              string       `json:"script"`
		Entrypoint          string       `json:"entrypoint"`
		User                string       `json:"user"`
		WorkspacePath       string       `json:"workspace"`
		Hostname            string       `json:"hostname"`
		Detach              bool         `json:"detach"`
		Healthcheck         *Healthcheck `json:"healthcheck"`
		Environment         []string     `json:"environment"`
		RequiredEnvironment []string     `json:"required_environment"`
	}

	Healthcheck struct {
		Command     string   `json:"command"`
		Interval    Duration `json:"interval"`
		Retries     int      `json:"retries"`
		StartPeriod Duration `json:"start_period"`
		Timeout     Duration `json:"timeout"`
	}

	Duration struct {
		time.Duration
	}
)

var ZeroDuration = Duration{time.Duration(0)}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	val := ""
	err := json.Unmarshal(data, &val)

	if err == nil {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration{parsed}
		return nil
	}

	return fmt.Errorf(
		"%s is not a valid duration",
		string(data),
	)
}
