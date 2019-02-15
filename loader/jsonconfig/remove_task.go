package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/util"
)

type RemoveTask struct {
	Extends             string          `json:"extends"`
	Environment         json.RawMessage `json:"environment"`
	RequiredEnvironment []string        `json:"required-environment"`
	Images              json.RawMessage `json:"images"`
	IncludeBuilt        bool            `json:"include-built"`
}

func (t *RemoveTask) Translate(name string) (config.Task, error) {
	images, err := util.UnmarshalStringList(t.Images)
	if err != nil {
		return nil, err
	}

	environment, err := util.UnmarshalStringList(t.Environment)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:                name,
		Extends:             t.Extends,
		Environment:         environment,
		RequiredEnvironment: t.RequiredEnvironment,
	}

	return &config.RemoveTask{
		TaskMeta:     meta,
		Images:       images,
		IncludeBuilt: t.IncludeBuilt,
	}, nil
}
