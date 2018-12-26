package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type PushTask struct {
	Extends     string          `json:"extends"`
	Images      json.RawMessage `json:"images"`
	Environment json.RawMessage `json:"environment"`
}

func (t *PushTask) Translate(name string) (config.Task, error) {
	images, err := unmarshalStringList(t.Images)
	if err != nil {
		return nil, err
	}

	environment, err := unmarshalStringList(t.Environment)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.PushTask{
		TaskMeta:    meta,
		Images:      images,
		Environment: environment,
	}, nil
}
