package config

import "fmt"

type RemoveTask struct {
	TaskMeta
	Images []string
}

func (t *RemoveTask) Extend(task Task) error {
	parent, ok := task.(*RemoveTask)
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
