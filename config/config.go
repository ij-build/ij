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

	// TODO - additional validation
	return nil
}

func (c *Config) validateTaskNames() error {
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
