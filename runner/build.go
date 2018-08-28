package runner

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type buildCommandBuilderState struct {
	state *state.State
	task  *config.BuildTask
	env   environment.Environment
}

func NewBuildCommandRunner(
	state *state.State,
	task *config.BuildTask,
	prefix *logging.Prefix,
	env environment.Environment,
) Runner {
	factory := buildCommandFactory(
		state,
		task,
		env,
	)

	return NewBaseRunner(
		state,
		prefix,
		NewMultiFactory(factory),
	)
}

func buildCommandFactory(
	state *state.State,
	task *config.BuildTask,
	env environment.Environment,
) BuilderFactory {
	return func() (*command.Builder, error) {
		s := &buildCommandBuilderState{
			state: state,
			task:  task,
			env:   env,
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

func (s *buildCommandBuilderState) addWorkspaceArg(cb *command.Builder) error {
	cb.AddArgs(s.state.Scratch.Workspace())
	return nil
}

func (s *buildCommandBuilderState) addDockerfileOptions(cb *command.Builder) error {
	dockerfile, err := s.env.ExpandString(s.task.Dockerfile)
	if err != nil {
		return err
	}

	cb.AddFlagValue("-f", dockerfile)
	return nil
}

func (s *buildCommandBuilderState) addTagOptions(cb *command.Builder) error {
	for _, tag := range s.task.Tags {
		expanded, err := s.env.ExpandString(tag)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-t", expanded)
	}

	return nil
}

func (s *buildCommandBuilderState) addLabelOptions(cb *command.Builder) error {
	for _, label := range s.task.Labels {
		expanded, err := s.env.ExpandString(label)
		if err != nil {
			return err
		}

		cb.AddFlagValue("--label", expanded)
	}

	return nil
}
