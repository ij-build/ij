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
		Hostname            string       `json:"hostname"`
		Detach              bool         `json:"detach"`
		Healthcheck         *Healthcheck `json:"healthcheck"`
		Environment         []string     `json:"environment"`
		RequiredEnvironment []string     `json:"required_environment"`
	}

	Healthcheck struct {
		Command     string   `json:"command"`
		Interval    duration `json:"interval"`
		Retries     int      `json:"retries"`
		StartPeriod duration `json:"start_period"`
		Timeout     duration `json:"timeout"`
	}

	duration struct {
		time.Duration
	}
)

var zeroDuration = duration{time.Duration(0)}

func (d duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *duration) UnmarshalJSON(data []byte) error {
	val := ""
	err := json.Unmarshal(data, &val)

	if err == nil {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = duration{parsed}
		return nil
	}

	return fmt.Errorf(
		"%s is not a valid duration",
		string(data),
	)
}
