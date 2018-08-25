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
	if _, ok := r.resolved[task.GetName()]; task.GetExtends() == "" || ok {
		return nil
	}

	if _, ok := r.resolving[task.GetName()]; ok {
		return fmt.Errorf(
			"failed to extend task %s (extension is cyclic)",
			task.GetName(),
		)
	}

	r.resolving[task.GetExtends()] = struct{}{}

	parent := r.config.Tasks[task.GetExtends()]
	if err := r.resolve(parent); err != nil {
		return err
	}

	if err := task.Extend(parent); err != nil {
		return err
	}

	r.resolved[task.GetExtends()] = struct{}{}
	return nil
}
