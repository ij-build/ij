package runner

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type PushTaskRunnerFactory func(
	*config.PushTask,
	environment.Environment,
	*logging.Prefix,
) TaskRunner

func NewPushTaskRunnerFactory(
	ctx context.Context,
	logger logging.Logger,
) PushTaskRunnerFactory {
	return func(
		task *config.PushTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		return NewBaseRunner(
			ctx,
			makePushTaskCommandFactory(task, env),
			logger,
			prefix,
		)
	}
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