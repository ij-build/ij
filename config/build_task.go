package config

import "fmt"

type BuildTask struct {
	TaskMeta
	Dockerfile  string
	Target      string
	Tags        []string
	Labels      []string
	Environment []string
}

func (t *BuildTask) GetType() string {
	return "build"
}

func (t *BuildTask) GetEnvironment() []string {
	return t.Environment
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
	t.Target = extendString(t.Target, parent.Target)
	t.Tags = append(parent.Tags, t.Tags...)
	t.Labels = append(parent.Labels, t.Labels...)
	t.Environment = append(parent.Environment, t.Environment...)
	return nil
}
