package config

import (
	"encoding/json"
	"fmt"
)

type RemoveTask struct {
	TaskMeta
	Images []string `json:"images,omitempty"`
}

func (t *RemoveTask) GetType() string {
	return "remove"
}

func (t *RemoveTask) Extend(task Task) error {
	parent, ok := task.(*RemoveTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.extendMeta(parent.TaskMeta)
	t.Images = append(parent.Images, t.Images...)
	return nil
}

func (t *RemoveTask) MarshalJSON() ([]byte, error) {
	type Alias RemoveTask

	return json.Marshal(&struct {
		*Alias
		Type string `json:"type"`
	}{
		Alias: (*Alias)(t),
		Type:  t.GetType(),
	})
}
