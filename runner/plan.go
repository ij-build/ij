package runner

import (
	"context"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
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

func (r *PlanRunner) IsDisabled(context *RunContext, name string) (bool, error) {
	if plan, ok := r.config.Plans[name]; ok {
		env := environment.Merge(
			environment.New(r.config.Environment),
			context.Environment,
			environment.New(plan.Environment),
			environment.New(context.GetExportedEnv()),
			environment.New(r.env),
		)

		val, err := env.ExpandString(plan.Disabled)
		if err != nil {
			return false, err
		}

		return val != "", nil
	}

	return false, nil
}

func (r *PlanRunner) ShouldRun(context *RunContext, name string) bool {
	if plans, ok := r.config.Metaplans[name]; ok {
		for _, plan := range plans {
			if r.ShouldRun(context, plan) {
				return true
			}
		}
	} else {
		for _, stage := range r.config.Plans[name].Stages {
			if stage.ShouldRun(context.Failure) && len(stage.Tasks) > 0 {
				return true
			}
		}
	}

	return false
}

func (r *PlanRunner) Run(
	context *RunContext,
	name string,
	prefix *logging.Prefix,
) bool {
	prefix = prefix.Append(name)

	disabled, err := r.IsDisabled(context, name)
	if err != nil {
		r.logger.Info(
			prefix,
			"Failed to expand environment for disabled check: %s",
			err.Error(),
		)

		return false
	}

	if disabled {
		r.logger.Warn(
			prefix,
			"Plan is disabled",
		)

		return true
	}

	if !r.ShouldRun(context, name) {
		r.logger.Info(
			prefix,
			"No tasks to perform",
		)

		return true
	}

	failure := context.Failure

	if plans, ok := r.config.Metaplans[name]; ok {
		r.logger.Info(
			prefix,
			"Beginning metaplan",
		)

		for _, plan := range plans {
			newContext := NewRunContext(context)
			newContext.Failure = failure

			if !r.Run(newContext, plan, prefix) {
				failure = true
			}
		}
	} else {
		r.logger.Info(
			prefix,
			"Beginning plan",
		)

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

		if !stage.ShouldRun(failure) || len(stage.Tasks) == 0 {
			r.logger.Info(
				stagePrefix,
				"No tasks to perform",
			)

			continue
		}

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
