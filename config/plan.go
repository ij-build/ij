package config

import "encoding/json"

type Plan struct {
	Name        string
	Stages      []*Stage `json:"stages"`
	Environment []string `json:"environment"`
}

type Stage struct {
	Name        string            `json:"name"`
	RawTasks    []json.RawMessage `json:"tasks"`
	Concurrent  bool              `json:"concurrent"`
	Environment []string          `json:"environment"`

	Tasks []*StageTask
}

type StageTask struct {
	Name        string   `json:"name"`
	Environment []string `json:"environment"`
}
