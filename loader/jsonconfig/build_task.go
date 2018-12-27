package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type BuildTask struct {
	Extends             string          `json:"extends"`
	Environment         json.RawMessage `json:"environment"`
	RequiredEnvironment []string        `json:"required_environment"`
	Dockerfile          string          `json:"dockerfile"`
	Target              string          `json:"target"`
	Tags                json.RawMessage `json:"tags"`
	Labels              json.RawMessage `json:"labels"`
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

	if t.Dockerfile == "" {
		t.Dockerfile = "Dockerfile"
	}

	return &config.BuildTask{
		TaskMeta:   meta,
		Dockerfile: t.Dockerfile,
		Target:     t.Target,
		Tags:       tags,
		Labels:     labels,
	}, nil
}
