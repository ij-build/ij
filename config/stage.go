package config

import "fmt"

type (
	Stage struct {
		Name        string       `json:"name,omitempty"`
		Disabled    string       `json:"disabled,omitempty"`
		BeforeStage string       `json:"before-stage,omitempty"`
		AfterStage  string       `json:"after-stage,omitempty"`
		RunMode     RunMode      `json:"run-mode,omitempty"`
		Parallel    bool         `json:"parallel,omitempty"`
		Environment []string     `json:"environment,omitempty"`
		Tasks       []*StageTask `json:"tasks,omitempty"`
	}

	StageTask struct {
		Name        string   `json:"name,omitempty"`
		Disabled    string   `json:"disabled,omitempty"`
		Environment []string `json:"environment,omitempty"`
	}

	RunMode int
)

const (
	_ RunMode = iota
	RunModeOnSuccess
	RunModeOnFailure
	RunModeAlways
)

func (s *Stage) ShouldRun(failure bool) bool {
	switch s.RunMode {
	case RunModeAlways:
		return true
	case RunModeOnSuccess:
		return !failure
	case RunModeOnFailure:
		return failure
	}

	return false
}

func (m RunMode) MarshalJSON() ([]byte, error) {
	switch m {
	case RunModeOnSuccess:
		return []byte(`"on-success"`), nil
	case RunModeOnFailure:
		return []byte(`"on-failure"`), nil
	case RunModeAlways:
		return []byte(`"always"`), nil
	}

	return nil, fmt.Errorf("unknown run-mode")
}
