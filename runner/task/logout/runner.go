package logout

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type Runner struct {
	state  *state.State
	task   *config.LogoutTask
	prefix *logging.Prefix
	env    environment.Environment
}

func NewRunner(
	state *state.State,
	task *config.LogoutTask,
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

	server, err := r.env.ExpandString(r.task.Server)
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	err = command.NewRunner(r.state.Logger).Run(
		r.state.Context,
		[]string{"docker", "logout", server},
		r.prefix,
	)

	if err != nil {
		r.state.ReportError(
			r.prefix,
			"Command failed: %s",
			err.Error(),
		)

		return false
	}

	return true
}
