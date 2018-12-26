package config

import "fmt"

type RemoveTask struct {
	TaskMeta
	Images      []string
	Environment []string
}

func (t *RemoveTask) GetType() string {
	return "remove"
}

func (t *RemoveTask) GetEnvironment() []string {
	return t.Environment
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

	t.Images = append(parent.Images, t.Images...)
	t.Environment = append(parent.Environment, t.Environment...)
	return nil
}
