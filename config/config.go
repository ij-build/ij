package config

import "fmt"

type Config struct {
	Extends     string
	Workspace   string
	Environment []string
	Tasks       map[string]Task
	Plans       map[string]*Plan
	Metaplans   map[string][]string
	Imports     []string
	Exports     []string
	Excludes    []string
	SSHIdentities []string
}

func (c *Config) Merge(child *Config) error {
	c.Environment = append(c.Environment, child.Environment...)

	for name, task := range child.Tasks {
		c.Tasks[name] = task
	}

	for name, plan := range child.Plans {
		if !plan.Extend {
			c.Plans[name] = plan
			continue
		}

		parentPlan, ok := c.Plans[name]
		if !ok {
			return fmt.Errorf(
				"plan %s extends unknown plan in parent",
				name,
			)
		}

		if err := parentPlan.Merge(plan); err != nil {
			return err
		}

		c.Plans[name] = parentPlan
	}

	return nil
}

func (c *Config) Resolve() error {
	resolver := NewTaskExtendsResolver(c)

	for _, task := range c.Tasks {
		if task.GetExtends() != "" {
			resolver.Add(task)
		}
	}

	return resolver.Resolve()
}

func (c *Config) Validate() error {
	validators := []func() error{
		c.validateTaskNames,
		c.validatePlanNames,
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) IsPlanDefined(name string) bool {
	if _, ok := c.Plans[name]; ok {
		return true
	}

	if _, ok := c.Metaplans[name]; ok {
		return true
	}

	return false
}

func (c *Config) validateTaskNames() error {
	for _, task := range c.Tasks {
		if _, ok := c.Tasks[task.GetExtends()]; task.GetExtends() != "" && !ok {
			return fmt.Errorf(
				"unknown task name %s referenced in task %s",
				task.GetExtends(),
				task.GetName(),
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

func (c *Config) validatePlanNames() error {
	for name, plans := range c.Metaplans {
		if _, ok := c.Plans[name]; ok {
			return fmt.Errorf(
				"plan is %s defined twice",
				name,
			)
		}

		for _, plan := range plans {
			if !c.IsPlanDefined(plan) {
				return fmt.Errorf(
					"unknown plan name %s referenced in metaplan %s",
					plan,
					name,
				)
			}
		}
	}

	return nil
}
