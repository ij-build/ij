package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
)

type (
	Stage struct {
		Name        string            `json:"name"`
		BeforeStage string            `json:"before_stage"`
		AfterStage  string            `json:"after_stage"`
		Tasks       []json.RawMessage `json:"tasks"`
		RunMode     string            `json:"run-mode"`
		Parallel    bool              `json:"parallel"`
		Environment []string          `json:"environment"`
	}

	StageTask struct {
		Name        string   `json:"name"`
		Environment []string `json:"environment"`
	}
)

func (s *Stage) Translate() (*config.Stage, error) {
	stageTasks := []*config.StageTask{}
	for _, stageTask := range s.Tasks {
		unmarshalled, err := unmarshalStageTask(stageTask)
		if err != nil {
			return nil, err
		}

		stageTasks = append(stageTasks, unmarshalled)
	}

	runMode, err := translateRunMode(s.RunMode)
	if err != nil {
		return nil, err
	}

	return &config.Stage{
		Name:        s.Name,
		BeforeStage: s.BeforeStage,
		AfterStage:  s.AfterStage,
		RunMode:     runMode,
		Parallel:    s.Parallel,
		Environment: s.Environment,
		Tasks:       stageTasks,
	}, nil
}

func translateRunMode(value string) (config.RunMode, error) {
	switch value {
	case "":
		fallthrough
	case "on-success":
		return config.RunModeOnSuccess, nil
	case "on-failure":
		return config.RunModeOnFailure, nil
	case "always":
		return config.RunModeAlways, nil
	}

	return 0, fmt.Errorf("unknown run mode '%s'", value)
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
