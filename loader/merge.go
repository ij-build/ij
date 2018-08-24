package loader

import (
	"fmt"

	"github.com/efritz/ij/config"
)

func mergeConfigs(child, parent *config.Config) error {
	parent.Environment = append(
		parent.Environment,
		child.Environment...,
	)

	for name, task := range child.Tasks {
		parent.Tasks[name] = task
	}

	for name, plan := range child.Plans {
		if !plan.Extend {
			parent.Plans[name] = plan
			continue
		}

		parentPlan, ok := parent.Plans[name]
		if !ok {
			return fmt.Errorf(
				"plan %s extends unknown plan in parent",
				name,
			)
		}

		if err := mergePlans(plan, parentPlan); err != nil {
			return err
		}

		parent.Plans[name] = parentPlan
	}

	return nil
}

func mergePlans(child, parent *config.Plan) error {
	for _, stage := range child.Stages {
		if err := mergeStage(parent, stage); err != nil {
			return err
		}
	}

	return nil
}

func mergeStage(parent *config.Plan, stage *config.Stage) error {
	if stage.BeforeStage != "" && stage.AfterStage != "" {
		return fmt.Errorf(
			"before_stage and after_stage declared in %s/%s",
			parent.Name,
			stage.Name,
		)
	}

	if index := stageIndex(parent, stage.Name); index >= 0 {
		if stage.BeforeStage != "" || stage.AfterStage != "" {
			return fmt.Errorf(
				"%s/%s exists in parent config, but before_stage or after_stage is also declared",
				parent.Name,
				stage.Name,
			)
		}

		parent.Stages[index] = stage
		return nil
	}

	targetStage := stage.Name

	if stage.BeforeStage != "" {
		if index := stageIndex(parent, stage.BeforeStage); index >= 0 {
			insertAtIndex(parent, stage, index)
			return nil
		}

		targetStage = stage.BeforeStage
	}

	if stage.AfterStage != "" {
		if index := stageIndex(parent, stage.AfterStage); index >= 0 {
			insertAtIndex(parent, stage, index+1)
			return nil
		}

		targetStage = stage.AfterStage
	}

	return fmt.Errorf(
		"stage %s/%s not declared in parent config",
		parent.Name,
		targetStage,
	)
}

func stageIndex(plan *config.Plan, name string) int {
	if name == "" {
		return -1
	}

	for i, stage := range plan.Stages {
		if stage.Name == name {
			return i
		}
	}

	return -1
}

func insertAtIndex(parent *config.Plan, stage *config.Stage, index int) {
	parent.Stages = append(parent.Stages, nil)
	copy(parent.Stages[index+1:], parent.Stages[index:])
	parent.Stages[index] = stage
}
