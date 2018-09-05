package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type BuildTask struct {
	Extends    string          `json:"extends"`
	Dockerfile string          `json:"dockerfile"`
	Tags       json.RawMessage `json:"tags"`
	Labels     json.RawMessage `json:"labels"`
}

func (t *BuildTask) Translate(name string) (config.Task, error) {
	tags, err := unmarshalStringList(t.Tags)
	if err != nil {
		return nil, err
	}

	labels, err := unmarshalStringList(t.Labels)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.BuildTask{
		TaskMeta:   meta,
		Dockerfile: t.Dockerfile,
		Tags:       tags,
		Labels:     labels,
	}, nil
}
