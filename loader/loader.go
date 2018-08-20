package loader

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/pvc/config"
)

func LoadPath(path string) (*config.Config, error) {
	data, err := readPath(path)
	if err != nil {
		return nil, err
	}

	return Load(data)
}

// TODO - load url and other stuff

func Load(data []byte) (*config.Config, error) {
	if err := validateWithSchema(data); err != nil {
		return nil, err
	}

	config := &config.Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if err := unmarshalStageTasks(config); err != nil {
		return nil, err
	}

	populateTaskNames(config)
	populatePlanNames(config)

	// TODO - additional validation
	if err := validateTaskNames(config); err != nil {
		return nil, err
	}

	return config, nil
}

func unmarshalStageTasks(config *config.Config) error {
	for _, plan := range config.Plans {
		for _, stage := range plan.Stages {
			for _, rawStageTask := range stage.RawTasks {
				stageTask, err := unmarshalStageTask(rawStageTask)
				if err != nil {
					return err
				}

				stage.Tasks = append(stage.Tasks, stageTask)
			}
		}
	}

	return nil
}

func unmarshalStageTask(raw json.RawMessage) (*config.StageTask, error) {
	stageTask := &config.StageTask{}
	if err := json.Unmarshal(raw, &stageTask.Name); err == nil {
		return stageTask, nil
	}

	if err := json.Unmarshal(raw, &stageTask); err != nil {
		return nil, err
	}

	return stageTask, nil
}

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

func validateTaskNames(config *config.Config) error {
	for _, plan := range config.Plans {
		for _, stage := range plan.Stages {
			for i, stageTask := range stage.Tasks {
				if _, ok := config.Tasks[stageTask.Name]; !ok {
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
