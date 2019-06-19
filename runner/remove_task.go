package runner

import (
	"context"

	"github.com/ij-build/ij/command"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
)

type RemoveTaskRunnerFactory func(
	*RunContext,
	*config.RemoveTask,
	environment.Environment,
	*logging.Prefix,
) TaskRunner

func NewRemoveTaskRunnerFactory(
	ctx context.Context,
	logger logging.Logger,
) RemoveTaskRunnerFactory {
	return func(
		context *RunContext,
		task *config.RemoveTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		return NewBaseRunner(
			ctx,
			removeTaskComandFactory(context, task, env),
			logger,
			prefix,
		)
	}
}

func removeTaskComandFactory(
	context *RunContext,
	task *config.RemoveTask,
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
				"rmi",
				"-f",
				image,
			}, nil)

			builders = append(builders, builder)
		}

		return builders, nil
	}
}
