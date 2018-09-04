package runner

import (
	"context"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/util"
)

type (
	StageRunner struct {
		ctx               context.Context
		logger            logging.Logger
		config            *config.Config
		taskRunnerFactory TaskRunnerFactory
		plan              *config.Plan
		stage             *config.Stage
		prefix            *logging.Prefix
		env               []string
	}

	TaskRunnerFunc func() bool
)

func NewStageRunner(
	ctx context.Context,
	logger logging.Logger,
	config *config.Config,
	taskRunnerFactory TaskRunnerFactory,
	plan *config.Plan,
	stage *config.Stage,
	prefix *logging.Prefix,
	env []string,
) *StageRunner {
	return &StageRunner{
		ctx:               ctx,
		logger:            logger,
		config:            config,
		taskRunnerFactory: taskRunnerFactory,
		plan:              plan,
		stage:             stage,
		prefix:            prefix,
		env:               env,
	}
}

func (r *StageRunner) Run(context *RunContext) bool {
	r.logger.Info(
		r.prefix,
		"Beginning stage",
	)

	var (
		runners   = []TaskRunnerFunc{}
		ambiguous = map[string]struct{}{}
		names     = map[string]struct{}{}
	)

	for _, stageTask := range r.stage.Tasks {
		if _, ok := names[stageTask.Name]; ok {
			ambiguous[stageTask.Name] = struct{}{}
		}

		names[stageTask.Name] = struct{}{}
	}

	for i, stageTask := range r.stage.Tasks {
		_, ok := ambiguous[stageTask.Name]

		runners = append(runners, r.buildTaskRunnerFunc(
			stageTask,
			i,
			r.config.Tasks[stageTask.Name],
			context,
			ok,
		))
	}

	if !r.stage.Parallel || r.config.Options.ForceSequential {
		return runSequential(runners)
	}

	return runParallel(runners)
}

func (r *StageRunner) buildTaskRunnerFunc(
	stageTask *config.StageTask,
	index int,
	task config.Task,
	context *RunContext,
	ambiguous bool,
) TaskRunnerFunc {
	name := task.GetName()

	if ambiguous {
		name = fmt.Sprintf(
			"%s.%d",
			name,
			index,
		)
	}

	taskPrefix := r.prefix.Append(name)

	return func() bool {
		env := environment.Merge(
			environment.New(r.config.Environment),
			environment.New(task.GetEnvironment()),
			context.Environment,
			environment.New(r.plan.Environment),
			environment.New(r.stage.Environment),
			environment.New(stageTask.Environment),
			environment.New(context.GetExportedEnv()),
			environment.New(r.env),
		)

		runner := r.taskRunnerFactory(
			task,
			taskPrefix,
			env,
		)

		if !runner.Run(context) {
			ReportError(
				r.ctx,
				r.logger,
				taskPrefix,
				"Task has failed",
			)

			return false
		}

		r.logger.Info(
			taskPrefix,
			"Task has completed successfully",
		)

		return true
	}
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
	for _, r := range runners {
		runner := r

		funcs = append(funcs, func() {
			if ok := runner(); !ok {
				failures <- false
			}
		})
	}

	util.RunParallel(funcs...).Wait()

	select {
	case <-failures:
		return false
	default:
	}

	return true
}
