package runner

import (
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

const PrefixMaxLength = 20

type (
	PlanTaskRunnerFactory func(
		*config.PlanTask,
		environment.Environment,
		*logging.Prefix,
	) TaskRunner

	planTaskRunner struct {
		runner *PlanRunner
		logger logging.Logger
		task   *config.PlanTask
		env    environment.Environment
		prefix *logging.Prefix
	}
)

func NewPlanTaskRunnerFactory(
	runner *PlanRunner,
	logger logging.Logger,
) PlanTaskRunnerFactory {
	return func(
		task *config.PlanTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		return &planTaskRunner{
			runner: runner,
			logger: logger,
			task:   task,
			env:    env,
			prefix: prefix,
		}
	}
}

func (r *planTaskRunner) Run(context *RunContext) bool {
	r.logger.Info(
		r.prefix,
		"Beginning task",
	)

	if r.prefix.Len() > PrefixMaxLength {
		r.logger.Error(
			r.prefix,
			"plan call history exceeds max depth",
		)

		return false
	}

	newContext := NewRunContext(context)
	newContext.Environment = r.env

	return r.runner.Run(newContext, r.task.Name, r.prefix)
}
