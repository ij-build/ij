package config

import "fmt"

type TaskExtendsResolver struct {
	config    *Config
	tasks     []*Task
	resolved  map[string]struct{}
	resolving map[string]struct{}
}

func NewTaskExtendsResolver(config *Config) *TaskExtendsResolver {
	return &TaskExtendsResolver{
		config:    config,
		tasks:     []*Task{},
		resolved:  map[string]struct{}{},
		resolving: map[string]struct{}{},
	}
}

func (r *TaskExtendsResolver) Add(task *Task) {
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

func (r *TaskExtendsResolver) resolve(task *Task) error {
	if _, ok := r.resolved[task.Name]; task.Extends == "" || ok {
		return nil
	}

	if _, ok := r.resolving[task.Name]; ok {
		return fmt.Errorf(
			"failed to extend task %s (extension is cyclic)",
			task.Name,
		)
	}

	r.resolving[task.Extends] = struct{}{}

	parent := r.config.Tasks[task.Extends]
	if err := r.resolve(parent); err != nil {
		return err
	}

	mergeTask(task, parent)
	r.resolved[task.Extends] = struct{}{}
	return nil
}

func mergeTask(child, parent *Task) {
	if child.Image == "" {
		child.Image = parent.Image
	}

	if child.Command == "" {
		child.Command = parent.Command
	}

	if child.Shell == "" {
		child.Shell = parent.Shell
	}

	if child.Script == "" {
		child.Script = parent.Script
	}

	if child.Entrypoint == "" {
		child.Entrypoint = parent.Entrypoint
	}

	if child.Hostname == "" {
		child.Hostname = parent.Hostname
	}

	if parent.Detach {
		child.Detach = true
	}

	child.Healthcheck = mergeHealthcheck(
		child.Healthcheck,
		parent.Healthcheck,
	)

	child.Environment = append(
		parent.Environment,
		child.Environment...,
	)

	child.RequiredEnvironment = append(
		parent.RequiredEnvironment,
		child.RequiredEnvironment...,
	)
}

func mergeHealthcheck(child, parent *Healthcheck) *Healthcheck {
	if child == nil {
		return parent
	}

	if parent == nil {
		return child
	}

	if child.Command == "" {
		child.Command = parent.Command
	}

	if child.Interval == ZeroDuration {
		child.Interval = parent.Interval
	}

	if child.Retries == 0 {
		child.Retries = parent.Retries
	}

	if child.StartPeriod == ZeroDuration {
		child.StartPeriod = parent.StartPeriod
	}

	if child.Timeout == ZeroDuration {
		child.Timeout = parent.Timeout
	}

	return child
}
