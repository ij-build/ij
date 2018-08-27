package config

import "fmt"

type BuildTask struct {
	TaskMeta
	Dockerfile string
	Tags       []string
	Labels     []string
	Arguments  []string
}

func (t *BuildTask) Extend(task Task) error {
	parent, ok := task.(*BuildTask)
	if !ok {
		return fmt.Errorf(
			"task %s extends %s, but they have different types",
			t.Name,
			task.GetName(),
		)
	}

	t.Dockerfile = extendString(t.Dockerfile, parent.Dockerfile)
	t.Tags = append(parent.Tags, t.Tags...)
	t.Labels = append(parent.Labels, t.Labels...)
	t.Arguments = append(parent.Arguments, t.Arguments...)
	return nil
}
