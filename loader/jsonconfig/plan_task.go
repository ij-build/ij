package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type PlanTask struct {
	Extends             string          `json:"extends"`
	Environment         json.RawMessage `json:"environment"`
	RequiredEnvironment []string        `json:"required_environment"`
	Name                string          `json:"name"`
}

func (t *PlanTask) Translate(name string) (config.Task, error) {
	environment, err := unmarshalStringList(t.Environment)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:                name,
		Extends:             t.Extends,
		Environment:         environment,
		RequiredEnvironment: t.RequiredEnvironment,
	}

	return &config.PlanTask{
		TaskMeta: meta,
		Name:     t.Name,
	}, nil
}
