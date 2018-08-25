package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type (
	Stage struct {
		Name        string            `json:"name"`
		BeforeStage string            `json:"before_stage"`
		AfterStage  string            `json:"after_stage"`
		Tasks       []json.RawMessage `json:"tasks"`
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

	return &config.Stage{
		Name:        s.Name,
		BeforeStage: s.BeforeStage,
		AfterStage:  s.AfterStage,
		Parallel:    s.Parallel,
		Environment: s.Environment,
		Tasks:       stageTasks,
	}, nil
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
