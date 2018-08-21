package config

import "encoding/json"

type (
	Plan struct {
		Name        string
		Stages      []*Stage `json:"stages"`
		Environment []string `json:"environment"`
	}

	Stage struct {
		Name        string            `json:"name"`
		RawTasks    []json.RawMessage `json:"tasks"`
		Parallel    bool              `json:"parallel"`
		Environment []string          `json:"environment"`

		Tasks []*StageTask
	}

	StageTask struct {
		Name        string   `json:"name"`
		Environment []string `json:"environment"`
	}
)
