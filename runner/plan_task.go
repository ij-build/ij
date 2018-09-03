package runner

import (
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

const PrefixMaxLength = 20

type planTaskRunner struct {
	state  *state.State
	task   *config.PlanTask
	prefix *logging.Prefix
	env    environment.Environment
}

func NewPlanTaskRunner(
	state *state.State,
	task *config.PlanTask,
	prefix *logging.Prefix,
	env environment.Environment,
) TaskRunner {
	return &planTaskRunner{
		state:  state,
		task:   task,
		prefix: prefix,
		env:    env,
	}
}

func (r *planTaskRunner) Run(context *RunContext) bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning task",
	)

	if r.prefix.Len() > PrefixMaxLength {
		r.state.Logger.Error(
			r.prefix,
			"plan call history exceeds max depth",
		)

		return false
	}

	return NewPlanRunner(r.state).Run(
		r.task.Name,
		r.prefix,
		&RunContext{
			Failure:     context.Failure,
			Environment: r.env,
		},
	)
}
