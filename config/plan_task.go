package config

import "fmt"

type PlanTask struct {
	TaskMeta
	Name string
}

func (t *PlanTask) GetType() string {
	return "plan"
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

	t.extendMeta(parent.TaskMeta)
	t.Name = extendString(t.Name, parent.Name)
	return nil
}
