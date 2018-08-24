package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Extends     string              `json:"extends"`
	Workspace   string              `json:"workspace"`
	Environment []string            `json:"environment"`
	Tasks       map[string]*Task    `json:"tasks"`
	Plans       map[string]*Plan    `json:"plans"`
	Metaplans   map[string][]string `json:"metaplans"`
	RawImports  json.RawMessage     `json:"import"`
	RawExports  json.RawMessage     `json:"export"`
	RawExcludes json.RawMessage     `json:"exclude"`

	Imports  []string
	Exports  []string
	Excludes []string
}

func (c *Config) Validate() error {
	if err := c.validateTaskNames(); err != nil {
		return err
	}

	if err := c.resolveTasks(); err != nil {
		return err
	}

	if err := c.validatePlanNames(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateTaskNames() error {
	for _, task := range c.Tasks {
		if _, ok := c.Tasks[task.Extends]; task.Extends != "" && !ok {
			return fmt.Errorf(
				"unknown task name %s referenced in task %s",
				task.Extends,
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

func (c *Config) IsPlanDefined(name string) bool {
	if _, ok := c.Plans[name]; ok {
		return true
	}

	if _, ok := c.Metaplans[name]; ok {
		return true
	}

	return false
}
