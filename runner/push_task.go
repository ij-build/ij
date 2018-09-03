package runner

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

func NewPushTaskRunner(
	state *state.State,
	task *config.PushTask,
	prefix *logging.Prefix,
	env environment.Environment,
) TaskRunner {
	return NewBaseRunner(
		state,
		prefix,
		makePushTaskCommandFactory(task, env),
	)
}

func makePushTaskCommandFactory(
	task *config.PushTask,
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
				"push",
				image,
			}, nil)

			builders = append(builders, builder)
		}

		return builders, nil
	}
}
