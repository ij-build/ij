package config

import "fmt"

type LogoutTask struct {
	TaskMeta
	Servers []string
}

func (t *LogoutTask) GetType() string {
	return "logout"
}

func (t *LogoutTask) Extend(task Task) error {
	parent, ok := task.(*LogoutTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.Servers = append(t.Servers, parent.Servers...)
	return nil
}
