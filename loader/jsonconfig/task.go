package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
)

type (
	TaskHint struct {
		Type string `json:"type"`
	}

	Task interface {
		Translate(name string) (config.Task, error)
	}
)

func translateTask(name string, data json.RawMessage) (config.Task, error) {
	hint := &TaskHint{Type: "run"}
	if err := json.Unmarshal(data, hint); err != nil {
		return nil, err
	}

	structMap := map[string]Task{
		"run":   &RunTask{},
		"build": &BuildTask{},
	}

	if task, ok := structMap[hint.Type]; ok {
		if err := json.Unmarshal(data, task); err != nil {
			return nil, err
		}

		return task.Translate(name)
	}

	return nil, fmt.Errorf("unknown task type '%s'", hint.Type)
}
