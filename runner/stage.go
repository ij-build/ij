package runner

import (
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
	"github.com/efritz/ij/task/run"
	"github.com/efritz/ij/util"
)

type (
	StageRunner struct {
		state  *state.State
		plan   *config.Plan
		stage  *config.Stage
		prefix *logging.Prefix
	}

	Runner interface {
		Run() bool
	}

	RunnerFunc func() bool
)

func NewStageRunner(
	state *state.State,
	plan *config.Plan,
	stage *config.Stage,
	prefix *logging.Prefix,
) *StageRunner {
	return &StageRunner{
		state:  state,
		plan:   plan,
		stage:  stage,
		prefix: prefix,
	}
}

func (r *StageRunner) Run() bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning stage",
	)

	runners := []RunnerFunc{}
	for i, stageTask := range r.stage.Tasks {
		runners = append(runners, r.buildRunnerFunc(
			stageTask,
			i,
			r.state.Config.Tasks[stageTask.Name],
		))
	}

	if !r.stage.Parallel || r.state.ForceSequential {
		return runSequential(runners)
	}

	return runParallel(runners)
}

func (r *StageRunner) buildRunnerFunc(
	stageTask *config.StageTask,
	index int,
	task config.Task,
) RunnerFunc {
	taskPrefix := r.prefix.Append(fmt.Sprintf(
		"%s.%d",
		task.GetName(),
		index,
	))

	return func() bool {
		if !r.buildRunner(task, taskPrefix, r.buildEnvironment(stageTask, task)).Run() {
			r.state.ReportError(
				taskPrefix,
				"Task has failed",
			)

			return false
		}

		r.state.Logger.Info(
			taskPrefix,
			"Task has completed successfully",
		)

		return true
	}
}

func (r *StageRunner) buildEnvironment(
	stageTask *config.StageTask,
	task config.Task,
) environment.Environment {
	return environment.Merge(
		environment.New(r.state.Config.Environment),
		environment.New(task.GetEnvironment()),
		environment.New(r.plan.Environment),
		environment.New(r.stage.Environment),
		environment.New(stageTask.Environment),
		environment.New(r.state.GetExportedEnv()),
		environment.New(r.state.Env),
	)
}

func (r *StageRunner) buildRunner(
	task config.Task,
	taskPrefix *logging.Prefix,
	env environment.Environment,
) Runner {
	switch t := task.(type) {
	case *config.RunTask:
		return run.NewRunner(r.state, t, taskPrefix, env)
	}

	panic("unexpected task type")
}

//
// Helpers

func runSequential(runners []RunnerFunc) bool {
	for _, runner := range runners {
		if !runner() {
			return false
		}
	}

	return true
}

func runParallel(runners []RunnerFunc) bool {
	failures := make(chan bool, len(runners))
	defer close(failures)

	funcs := []func(){}
	for _, runner := range runners {
		funcs = append(funcs, func() {
			if ok := runner(); !ok {
				failures <- false
			}
		})
	}

	util.RunParallel(funcs...)

	select {
	case <-failures:
		return false
	default:
	}

	return true
}
