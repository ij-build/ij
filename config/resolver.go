package config

import "fmt"

type TaskExtendsResolver struct {
	config    *Config
	tasks     []Task
	resolved  map[string]struct{}
	resolving map[string]struct{}
}

func NewTaskExtendsResolver(config *Config) *TaskExtendsResolver {
	return &TaskExtendsResolver{
		config:    config,
		tasks:     []Task{},
		resolved:  map[string]struct{}{},
		resolving: map[string]struct{}{},
	}
}

func (r *TaskExtendsResolver) Add(task Task) {
	r.tasks = append(r.tasks, task)
}

func (r *TaskExtendsResolver) Resolve() error {
	for _, task := range r.tasks {
		if err := r.resolve(task); err != nil {
			return err
		}
	}

	return nil
}

func (r *TaskExtendsResolver) resolve(task Task) error {
	var (
		name    = task.GetName()
		extends = task.GetExtends()
	)

	if extends == "" {
		return nil
	}

	parent, ok := r.config.Tasks[extends]
	if !ok {
		return fmt.Errorf(
			"unknown task name %s referenced in task %s",
			extends,
			name,
		)
	}

	if _, ok := r.resolved[name]; ok {
		return nil
	}

	if _, ok := r.resolving[name]; ok {
		return fmt.Errorf(
			"failed to extend task %s (extension is cyclic)",
			name,
		)
	}

	r.resolving[name] = struct{}{}

	if err := r.resolve(parent); err != nil {
		return err
	}

	if err := task.Extend(parent); err != nil {
		return err
	}

	delete(r.resolving, name)
	r.resolved[name] = struct{}{}
	return nil
}
