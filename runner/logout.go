package runner

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

func NewLogoutCommandRunner(
	state *state.State,
	task *config.LogoutTask,
	prefix *logging.Prefix,
	env environment.Environment,
) Runner {
	return NewBaseRunner(
		state,
		prefix,
		logoutCommandFactory(task, env),
	)
}

func logoutCommandFactory(
	task *config.LogoutTask,
	env environment.Environment,
) BuilderSetFactory {
	return func() ([]*command.Builder, error) {
		servers, err := env.ExpandSlice(task.Servers)
		if err != nil {
			return nil, err
		}

		builders := []*command.Builder{}
		for _, server := range servers {
			builder := command.NewBuilder([]string{
				"docker",
				"logout",
				server,
			}, nil)

			builders = append(builders, builder)
		}

		return builders, nil
	}
}
