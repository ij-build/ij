package config

import (
	"encoding/json"
	"fmt"
)

type PushTask struct {
	TaskMeta
	Images       []string `json:"images,omitempty"`
	IncludeBuilt bool     `json:"include-built,omitempty"`
}

func (t *PushTask) GetType() string {
	return "push"
}

func (t *PushTask) Extend(task Task) error {
	parent, ok := task.(*PushTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.extendMeta(parent.TaskMeta)
	t.Images = append(parent.Images, t.Images...)
	t.IncludeBuilt = extendBool(t.IncludeBuilt, parent.IncludeBuilt)
	return nil
}

func (t *PushTask) MarshalJSON() ([]byte, error) {
	type Alias PushTask

	return json.Marshal(&struct {
		*Alias
		Type string `json:"type"`
	}{
		Alias: (*Alias)(t),
		Type:  t.GetType(),
	})
}
