package config

import "fmt"

type (
	Config struct {
		Extends       string
		Registries    []Registry
		SSHIdentities []string
		Environment   []string
		Import        *FileList
		Export        *FileList
		Workspace     string
		Tasks         map[string]Task
		Plans         map[string]*Plan
		Metaplans     map[string][]string
	}

	FileList struct {
		Files    []string
		Excludes []string
	}

	Override struct {
		Registries     []Registry
		SSHIdentities  []string
		Environment    []string
		ImportExcludes []string
		ExportExcludes []string
	}
)

func (c *Config) Merge(child *Config) error {
	c.Registries = append(c.Registries, child.Registries...)
	c.SSHIdentities = append(c.SSHIdentities, child.SSHIdentities...)
	c.Environment = append(c.Environment, child.Environment...)
	c.Import.Merge(child.Import)
	c.Export.Merge(child.Export)

	if child.Workspace != "" {
		c.Workspace = child.Workspace
	}

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

	for name, plans := range child.Metaplans {
		c.Metaplans[name] = plans
	}

	return nil
}

func (f *FileList) Merge(child *FileList) {
	f.Files = append(f.Files, child.Files...)
	f.Excludes = append(f.Excludes, child.Excludes...)
}

func (c *Config) ApplyOverride(override *Override) {
	c.Registries = append(c.Registries, override.Registries...)
	c.SSHIdentities = append(c.SSHIdentities, override.SSHIdentities...)
	c.Environment = append(c.Environment, override.Environment...)
	c.Import.Excludes = append(c.Import.Excludes, override.ImportExcludes...)
	c.Export.Excludes = append(c.Export.Excludes, override.ExportExcludes...)
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
			for _, stageTask := range stage.Tasks {
				if _, ok := c.Tasks[stageTask.Name]; !ok {
					return fmt.Errorf(
						"unknown task name %s referenced in %s/%s",
						stageTask.Name,
						plan.Name,
						stage.Name,
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
				"plan %s is defined twice",
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
