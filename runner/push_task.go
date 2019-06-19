package runner

import (
	"context"

	"github.com/ij-build/ij/command"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
)

type PushTaskRunnerFactory func(
	*RunContext,
	*config.PushTask,
	environment.Environment,
	*logging.Prefix,
) TaskRunner

func NewPushTaskRunnerFactory(
	ctx context.Context,
	logger logging.Logger,
) PushTaskRunnerFactory {
	return func(
		context *RunContext,
		task *config.PushTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		return NewBaseRunner(
			ctx,
			makePushTaskCommandFactory(context, task, env),
			logger,
			prefix,
		)
	}
}

func makePushTaskCommandFactory(
	context *RunContext,
	task *config.PushTask,
	env environment.Environment,
) BuilderSetFactory {
	return func() ([]*command.Builder, error) {
		images, err := env.ExpandSlice(task.Images)
		if err != nil {
			return nil, err
		}

		if task.IncludeBuilt {
			images = append(images, context.GetTags()...)
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
