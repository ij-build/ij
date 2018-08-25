package config

type (
	Stage struct {
		Name        string
		BeforeStage string
		AfterStage  string
		Parallel    bool
		Environment []string
		Tasks       []*StageTask
	}

	StageTask struct {
		Name        string
		Environment []string
	}
)
