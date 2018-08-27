package remove

import (
	"strings"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type Runner struct {
	state  *state.State
	task   *config.RemoveTask
	prefix *logging.Prefix
	env    environment.Environment
}

func NewRunner(
	state *state.State,
	task *config.RemoveTask,
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

	images, err := r.env.ExpandSlice(r.task.Images)
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	for _, image := range images {
		_, stderr, err := command.NewRunner(r.state.Logger).RunForOutput(
			r.state.Context,
			[]string{"docker", "rmi", image},
		)

		if err != nil {
			if strings.Contains(stderr, "No such image") {
				continue
			}

			r.state.ReportError(
				r.prefix,
				"Command failed: %s",
				err.Error(),
			)

			return false
		}
	}

	return true
}
