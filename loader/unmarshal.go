package loader

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

func unmarshalFileList(config *config.Config) error {
	imports, err := unmarshalStringList(config.RawImports)
	if err != nil {
		return err
	}

	exports, err := unmarshalStringList(config.RawExports)
	if err != nil {
		return err
	}

	excludes, err := unmarshalStringList(config.RawExcludes)
	if err != nil {
		return err
	}

	config.Imports = imports
	config.Exports = exports
	config.Excludes = excludes
	return nil
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

func unmarshalStringList(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	single := ""
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{single}, nil
	}

	multiple := []string{}
	if err := json.Unmarshal(raw, &multiple); err != nil {
		return nil, err
	}

	return multiple, nil
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
