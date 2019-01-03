package config

import (
	"fmt"
	"time"
)

type (
	Config struct {
		Extends          []string
		Options          *Options
		Registries       []Registry
		Workspace        string
		Environment      []string
		EnvironmentFiles []string
		Import           *ImportFileList
		Export           *ExportFileList
		Tasks            map[string]Task
		Plans            map[string]*Plan
		Metaplans        map[string][]string
	}

	Options struct {
		SSHIdentities       []string
		ForceSequential     bool
		HealthcheckInterval time.Duration
	}

	ImportFileList struct {
		Files    []string
		Excludes []string
	}

	ExportFileList struct {
		Files         []string
		Excludes      []string
		CleanExcludes []string
	}

	Override struct {
		Options          *Options
		Registries       []Registry
		Environment      []string
		EnvironmentFiles []string
		ImportExcludes   []string
		ExportExcludes   []string
		CleanExcludes    []string
	}
)

func (c *Config) Merge(child *Config) error {
	c.Options.Merge(child.Options)
	c.Registries = append(c.Registries, child.Registries...)
	c.Environment = append(c.Environment, child.Environment...)
	c.EnvironmentFiles = append(c.EnvironmentFiles, child.EnvironmentFiles...)
	c.Import.Merge(child.Import)
	c.Export.Merge(child.Export)

	c.Workspace = extendString(child.Workspace, c.Workspace)

	for name, task := range child.Tasks {
		c.Tasks[name] = task
	}

	for name, plan := range child.Plans {
		if plan.Extends == "" {
			c.Plans[name] = plan
			continue
		}

		parentPlan, ok := c.Plans[plan.Extends]
		if !ok {
			return fmt.Errorf(
				"plan %s extends unknown plan %s in parent",
				name,
				plan.Extends,
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

func (o *Options) Merge(child *Options) {
	if len(child.SSHIdentities) > 0 {
		o.SSHIdentities = child.SSHIdentities
	}

	o.ForceSequential = extendBool(child.ForceSequential, o.ForceSequential)
	o.HealthcheckInterval = extendDuration(child.HealthcheckInterval, o.HealthcheckInterval)
}

func (f *ImportFileList) Merge(child *ImportFileList) {
	f.Files = append(f.Files, child.Files...)
	f.Excludes = append(f.Excludes, child.Excludes...)
}

func (f *ExportFileList) Merge(child *ExportFileList) {
	f.Files = append(f.Files, child.Files...)
	f.Excludes = append(f.Excludes, child.Excludes...)
	f.CleanExcludes = append(f.CleanExcludes, child.CleanExcludes...)
}

func (c *Config) ApplyOverride(override *Override) {
	c.Options.Merge(override.Options)
	c.Registries = append(c.Registries, override.Registries...)
	c.Environment = append(c.Environment, override.Environment...)
	c.EnvironmentFiles = append(c.EnvironmentFiles, override.EnvironmentFiles...)
	c.Import.Excludes = append(c.Import.Excludes, override.ImportExcludes...)
	c.Export.Excludes = append(c.Export.Excludes, override.ExportExcludes...)
	c.Export.CleanExcludes = append(c.Export.CleanExcludes, override.CleanExcludes...)
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
	for name, task := range c.Tasks {
		if planTask, ok := task.(*PlanTask); ok {
			if !c.IsPlanDefined(planTask.Name) {
				return fmt.Errorf(
					"unknown plan name %s referenced in task %s",
					planTask.Name,
					name,
				)
			}
		}
	}

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
