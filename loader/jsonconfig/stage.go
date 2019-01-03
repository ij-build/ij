package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/util"
)

type (
	Stage struct {
		Name        string            `json:"name"`
		Disabled    string            `json:"disabled"`
		BeforeStage string            `json:"before-stage"`
		AfterStage  string            `json:"after-stage"`
		Tasks       []json.RawMessage `json:"tasks"`
		RunMode     string            `json:"run-mode"`
		Parallel    bool              `json:"parallel"`
		Environment json.RawMessage   `json:"environment"`
	}

	StageTask struct {
		Name        string          `json:"name"`
		Disabled    string          `json:"disabled"`
		Environment json.RawMessage `json:"environment"`
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

	environment, err := util.UnmarshalStringList(s.Environment)
	if err != nil {
		return nil, err
	}

	return &config.Stage{
		Name:        s.Name,
		Disabled:    s.Disabled,
		BeforeStage: s.BeforeStage,
		AfterStage:  s.AfterStage,
		Tasks:       stageTasks,
		RunMode:     runMode,
		Parallel:    s.Parallel,
		Environment: environment,
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
	var name string
	if err := json.Unmarshal(raw, &name); err == nil {
		return &config.StageTask{Name: name}, nil
	}

	stageTask := &StageTask{}
	if err := json.Unmarshal(raw, &stageTask); err != nil {
		return nil, err
	}

	environment, err := util.UnmarshalStringList(stageTask.Environment)
	if err != nil {
		return nil, err
	}

	return &config.StageTask{
		Name:        stageTask.Name,
		Disabled:    stageTask.Disabled,
		Environment: environment,
	}, nil
}
