package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader/schema"
)

type (
	TaskExtendHint struct {
		Extend string `json:"extends"`
	}

	TaskTypeHint struct {
		Type string `json:"type"`
	}

	Task interface {
		Translate(name string) (config.Task, error)
	}
)

func translateTask(
	parent *config.Config,
	name string,
	data json.RawMessage,
) (config.Task, error) {
	typeHint := &TaskTypeHint{Type: "run"}

	if parent != nil {
		extendHint := &TaskExtendHint{}
		if err := json.Unmarshal(data, extendHint); err != nil {
			return nil, err
		}

		if parentTask, ok := parent.Tasks[extendHint.Extend]; ok {
			typeHint.Type = parentTask.GetType()
		}
	}

	if err := json.Unmarshal(data, typeHint); err != nil {
		return nil, err
	}

	structMap := map[string]Task{
		"run":    &RunTask{},
		"build":  &BuildTask{},
		"push":   &PushTask{},
		"remove": &RemoveTask{},
	}

	task, ok := structMap[typeHint.Type]
	if !ok {
		return nil, fmt.Errorf("unknown task type '%s'", typeHint.Type)
	}

	assetName := fmt.Sprintf("schema/%s.yaml", typeHint.Type)

	if err := schema.Validate(assetName, data); err != nil {
		return nil, fmt.Errorf("failed to validate task %s: %s", name, err.Error())
	}

	if err := json.Unmarshal(data, task); err != nil {
		return nil, err
	}

	return task.Translate(name)
}
