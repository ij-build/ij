package runner

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

func NewRemoveCommandRunner(
	state *state.State,
	task *config.RemoveTask,
	prefix *logging.Prefix,
	env environment.Environment,
) Runner {
	return NewBaseRunner(
		state,
		prefix,
		removeComandFactory(task, env),
	)
}

func removeComandFactory(
	task *config.RemoveTask,
	env environment.Environment,
) BuilderSetFactory {
	return func() ([]*command.Builder, error) {
		images, err := env.ExpandSlice(task.Images)
		if err != nil {
			return nil, err
		}

		builders := []*command.Builder{}
		for _, image := range images {
			builder := command.NewBuilder([]string{
				"docker",
				"remove",
				image,
			}, nil)

			builders = append(builders, builder)
		}

		return builders, nil
	}
}
