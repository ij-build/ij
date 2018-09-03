package config

import "fmt"

type PlanTask struct {
	TaskMeta
	Name        string
	Environment []string
}

func (t *PlanTask) GetType() string {
	return "plan"
}

func (t *PlanTask) GetEnvironment() []string {
	return t.Environment
}

func (t *PlanTask) Extend(task Task) error {
	parent, ok := task.(*PlanTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.Name = extendString(t.Name, parent.Name)
	t.Environment = append(parent.Environment, t.Environment...)
	return nil
}
