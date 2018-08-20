package loader

import (
	"encoding/json"

	"github.com/efritz/pvc/config"
)

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
