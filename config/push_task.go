package config

import "fmt"

type PushTask struct {
	TaskMeta
	Images      []string
	Environment []string
}

func (t *PushTask) GetType() string {
	return "push"
}

func (t *PushTask) GetEnvironment() []string {
	return t.Environment
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

	t.Images = append(parent.Images, t.Images...)
	t.Environment = append(parent.Environment, t.Environment...)
	return nil
}
