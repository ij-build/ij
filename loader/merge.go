package loader

import "github.com/efritz/ij/config"

func merge(child, parent *config.Config) error {
	parent.Environment = append(
		parent.Environment,
		child.Environment...,
	)

	for name, task := range child.Tasks {
		parent.Tasks[name] = task
	}

	for name, plan := range child.Plans {
		parent.Plans[name] = plan
	}

	return nil
}
