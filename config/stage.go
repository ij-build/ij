package config

type (
	Stage struct {
		Name        string
		BeforeStage string
		AfterStage  string
		RunMode     RunMode
		Parallel    bool
		Environment []string
		Tasks       []*StageTask
	}

	StageTask struct {
		Name        string
		Environment []string
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
