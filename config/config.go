package config

import "fmt"

type Config struct {
	Extends     string           `json:"extends"`
	Environment []string         `json:"environment"`
	Tasks       map[string]*Task `json:"tasks"`
	Plans       map[string]*Plan `json:"plans"`
}

func (c *Config) Validate() error {
	if err := c.validateTaskNames(); err != nil {
		return err
	}

	if err := c.resolveTasks(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateTaskNames() error {
	for _, task := range c.Tasks {
		if _, ok := c.Tasks[task.Extends]; task.Extends != "" && !ok {
			return fmt.Errorf(
				"unknown task name %s referenced in task %s",
				task.Name,
			)
		}
	}

	for _, plan := range c.Plans {
		for _, stage := range plan.Stages {
			for i, stageTask := range stage.Tasks {
				if _, ok := c.Tasks[stageTask.Name]; !ok {
					return fmt.Errorf(
						"unknown task name %s referenced in %s/%s/%s #(%d)",
						stageTask.Name,
						plan.Name,
						stage.Name,
						stageTask.Name,
						i,
					)
				}
			}
		}
	}

	return nil
}

func (c *Config) resolveTasks() error {
	resolver := NewTaskExtendsResolver(c)

	for _, task := range c.Tasks {
		if task.Extends != "" {
			resolver.Add(task)
		}
	}

	if err := resolver.Resolve(); err != nil {
		return err
	}

	for _, task := range c.Tasks {
		if task.Image == "" {
			return fmt.Errorf(
				"no image supplied for task %s",
				task.Name,
			)
		}
	}

	return nil
}
