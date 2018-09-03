package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type PushTask struct {
	Extends string          `json:"extends"`
	Images  json.RawMessage `json:"images"`
}

func (t *PushTask) Translate(name string) (config.Task, error) {
	images, err := unmarshalStringList(t.Images)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.PushTask{
		TaskMeta: meta,
		Images:   images,
	}, nil
}
