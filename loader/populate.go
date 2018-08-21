package loader

import "github.com/efritz/ij/config"

func populateTaskNames(config *config.Config) {
	for name, task := range config.Tasks {
		task.Name = name
	}
}

func populatePlanNames(config *config.Config) {
	for name, plan := range config.Plans {
		plan.Name = name
	}
}
