package build

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type Runner struct {
	state  *state.State
	task   *config.BuildTask
	prefix *logging.Prefix
	env    environment.Environment
}

func NewRunner(
	state *state.State,
	task *config.BuildTask,
	prefix *logging.Prefix,
	env environment.Environment,
) *Runner {
	return &Runner{
		state:  state,
		task:   task,
		prefix: prefix,
		env:    env,
	}
}

func (r *Runner) Run() bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning task",
	)

	args, err := Build(r.state, r.task, r.env)
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	commandErr := command.Run(
		r.state.Context,
		args,
		r.state.Logger,
		r.prefix,
	)

	if commandErr != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			commandErr.Error(),
		)

		return false
	}

	return true
}
