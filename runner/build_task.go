package runner

import (
	"context"

	"github.com/ij-build/ij/command"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
)

type (
	BuildTaskRunnerFactory func(
		*config.BuildTask,
		environment.Environment,
		*logging.Prefix,
	) TaskRunner

	buildTaskCommandBuilderState struct {
		ctx       context.Context
		logger    logging.Logger
		workspace string
		env       environment.Environment
		task      *config.BuildTask
	}
)

func NewBuildTaskRunnerFactory(
	ctx context.Context,
	workspace string,
	logger logging.Logger,
) BuildTaskRunnerFactory {
	return func(
		task *config.BuildTask,
		env environment.Environment,
		prefix *logging.Prefix,
	) TaskRunner {
		factory := buildTaskCommandFactory(
			workspace,
			task,
			env,
		)

		runner := NewBaseRunner(
			ctx,
			NewMultiFactory(factory),
			logger,
			prefix,
		)

		runner.RegisterOnSuccess(func(context *RunContext) error {
			tags := []string{}
			for _, tag := range task.Tags {
				expanded, err := env.ExpandString(tag)
				if err != nil {
					return err
				}

				tags = append(tags, expanded)
			}

			context.AddTags(tags)
			return nil
		})

		return runner
	}
}

func buildTaskCommandFactory(
	workspace string,
	task *config.BuildTask,
	env environment.Environment,
) BuilderFactory {
	return func() (*command.Builder, error) {
		s := &buildTaskCommandBuilderState{
			workspace: workspace,
			task:      task,
			env:       env,
		}

		return command.NewBuilder(
			[]string{
				"docker",
				"build",
			},
			[]command.BuildFunc{
				s.addWorkspaceArg,
				s.addDockerfileOptions,
				s.addTargetOptions,
				s.addTagOptions,
				s.addLabelOptions,
			},
		), nil
	}
}

func (s *buildTaskCommandBuilderState) addWorkspaceArg(cb *command.Builder) error {
	cb.AddArgs(s.workspace)
	return nil
}

func (s *buildTaskCommandBuilderState) addDockerfileOptions(cb *command.Builder) error {
	dockerfile, err := s.env.ExpandString(s.task.Dockerfile)
	if err != nil {
		return err
	}

	cb.AddFlagValue("-f", dockerfile)
	return nil
}

func (s *buildTaskCommandBuilderState) addTargetOptions(cb *command.Builder) error {
	target, err := s.env.ExpandString(s.task.Target)
	if err != nil {
		return err
	}

	cb.AddFlagValue("--target", target)
	return nil
}

func (s *buildTaskCommandBuilderState) addTagOptions(cb *command.Builder) error {
	for _, tag := range s.task.Tags {
		expanded, err := s.env.ExpandString(tag)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-t", expanded)
	}

	return nil
}

func (s *buildTaskCommandBuilderState) addLabelOptions(cb *command.Builder) error {
	for _, label := range s.task.Labels {
		expanded, err := s.env.ExpandString(label)
		if err != nil {
			return err
		}

		cb.AddFlagValue("--label", expanded)
	}

	return nil
}
