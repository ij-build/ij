package runner

import (
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
	"github.com/efritz/ij/util"
)

type (
	StageRunner struct {
		state  *state.State
		plan   *config.Plan
		stage  *config.Stage
		prefix *logging.Prefix
	}

	TaskRunnerFunc func() bool
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

func (r *StageRunner) Run(context *RunContext) bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning stage",
	)

	runners := []TaskRunnerFunc{}
	for i, stageTask := range r.stage.Tasks {
		runners = append(runners, r.buildTaskRunnerFunc(
			stageTask,
			i,
			r.state.Config.Tasks[stageTask.Name],
			context,
		))
	}

	if !r.stage.Parallel || r.state.Config.Options.ForceSequential {
		return runSequential(runners)
	}

	return runParallel(runners)
}

func (r *StageRunner) buildTaskRunnerFunc(
	stageTask *config.StageTask,
	index int,
	task config.Task,
	context *RunContext,
) TaskRunnerFunc {
	taskPrefix := r.prefix.Append(fmt.Sprintf(
		"%s.%d",
		task.GetName(),
		index,
	))

	return func() bool {
		runner := r.buildRunner(
			task,
			taskPrefix,
			r.buildEnvironment(stageTask, context, task),
		)

		if !runner.Run(context) {
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
	context *RunContext,
	task config.Task,
) environment.Environment {
	return r.state.BuildEnv(
		environment.New(task.GetEnvironment()),
		context.Environment,
		environment.New(r.plan.Environment),
		environment.New(r.stage.Environment),
		environment.New(stageTask.Environment),
	)
}

func (r *StageRunner) buildRunner(
	task config.Task,
	taskPrefix *logging.Prefix,
	env environment.Environment,
) TaskRunner {
	switch t := task.(type) {
	case *config.BuildTask:
		return NewBuildTaskRunner(r.state, t, taskPrefix, env)
	case *config.PushTask:
		return NewPushTaskRunner(r.state, t, taskPrefix, env)
	case *config.RemoveTask:
		return NewRemoveTaskRunner(r.state, t, taskPrefix, env)
	case *config.RunTask:
		return NewRunTaskRunner(r.state, t, taskPrefix, env)
	case *config.PlanTask:
		return NewPlanTaskRunner(r.state, t, taskPrefix, env)
	}

	panic("unexpected task type")
}

//
// Helpers

func runSequential(runners []TaskRunnerFunc) bool {
	for _, runner := range runners {
		if !runner() {
			return false
		}
	}

	return true
}

func runParallel(runners []TaskRunnerFunc) bool {
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
