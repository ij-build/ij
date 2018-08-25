package config

import "fmt"

type PushTask struct {
	TaskMeta
	Images []string
}

func (t *PushTask) Extend(task Task) error {
	parent, ok := task.(*PushTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			parent.GetName(),
		)
	}

	t.Images = append(parent.Images, t.Images...)
	return nil
}
