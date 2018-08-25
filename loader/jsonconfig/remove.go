package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type RemoveTask struct {
	Extends string          `json:"extends"`
	Images  json.RawMessage `json:"images"`
}

func (t *RemoveTask) Translate(name string) (config.Task, error) {
	images, err := unmarshalStringList(t.Images)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.RemoveTask{
		TaskMeta: meta,
		Images:   images,
	}, nil
}
