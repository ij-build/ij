package runner

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
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

		return NewBaseRunner(
			ctx,
			NewMultiFactory(factory),
			logger,
			prefix,
		)
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
