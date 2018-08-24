package config

import "encoding/json"

type (
	Plan struct {
		Name        string
		Extend      bool     `json:"extend"`
		Stages      []*Stage `json:"stages"`
		Environment []string `json:"environment"`
	}

	Stage struct {
		Name        string            `json:"name"`
		BeforeStage string            `json:"before_stage"`
		AfterStage  string            `json:"after_stage"`
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
