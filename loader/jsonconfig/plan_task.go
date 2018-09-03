package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type PlanTask struct {
	Extends     string          `json:"extends"`
	Name        string          `json:"name"`
	Environment json.RawMessage `json:"environment"`
}

func (t *PlanTask) Translate(name string) (config.Task, error) {
	environment, err := unmarshalStringList(t.Environment)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.PlanTask{
		TaskMeta:    meta,
		Name:        t.Name,
		Environment: environment,
	}, nil
}
