package runner

import (
	"context"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/logging"
)

type PlanRunner struct {
	ctx               context.Context
	config            *config.Config
	taskRunnerFactory TaskRunnerFactory
	logger            logging.Logger
	env               []string
}

func NewPlanRunner(
	ctx context.Context,
	config *config.Config,
	taskRunnerFactory TaskRunnerFactory,
	logger logging.Logger,
	env []string,
) *PlanRunner {
	return &PlanRunner{
		ctx:               ctx,
		config:            config,
		taskRunnerFactory: taskRunnerFactory,
		logger:            logger,
		env:               env,
	}
}

func (r *PlanRunner) Run(
	context *RunContext,
	name string,
	prefix *logging.Prefix,
) bool {
	prefix = prefix.Append(name)

	r.logger.Info(
		prefix,
		"Beginning plan",
	)

	failure := context.Failure

	if plans, ok := r.config.Metaplans[name]; ok {
		for _, plan := range plans {
			newContext := NewRunContext(context)
			newContext.Failure = failure

			if !r.Run(newContext, plan, prefix) {
				failure = true
			}
		}
	} else {
		newContext := NewRunContext(context)
		newContext.Failure = failure

		failure = !r.runPlan(newContext, name, prefix)
	}

	if failure {
		suffix := ""
		if context.Failure {
			suffix = " (due to previous failure)"
		}

		r.logger.Error(
			prefix,
			"Plan failed%s",
			suffix,
		)
	} else {
		r.logger.Info(
			prefix,
			"Plan completed successfully",
		)
	}

	return !failure
}

func (r *PlanRunner) runPlan(
	context *RunContext,
	name string,
	prefix *logging.Prefix,
) bool {
	var (
		plan    = r.config.Plans[name]
		failure = context.Failure
	)

	for _, stage := range plan.Stages {
		stagePrefix := prefix.Append(stage.Name)

		if !stage.ShouldRun(context.Failure) {
			r.logger.Info(
				stagePrefix,
				"Skipping stage",
			)

			continue
		}

		runner := NewStageRunner(
			r.ctx,
			r.logger,
			r.config,
			r.taskRunnerFactory,
			plan,
			stage,
			stagePrefix,
			r.env,
		)

		newContext := NewRunContext(context)
		newContext.Failure = failure

		if !runner.Run(newContext) {
			r.logger.Error(
				stagePrefix,
				"Stage failed",
			)

			failure = true
			continue
		}

		r.logger.Info(
			stagePrefix,
			"Stage completed successfully",
		)
	}

	return !failure
}
