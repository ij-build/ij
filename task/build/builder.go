package build

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/state"
)

type builder struct {
	state *state.State
	task  *config.BuildTask
	env   environment.Environment
}

func Build(
	state *state.State,
	task *config.BuildTask,
	env environment.Environment,
) ([]string, error) {
	b := &builder{
		state: state,
		task:  task,
		env:   env,
	}

	builderFuncs := []command.BuildFunc{
		b.addDockerfileOptions,
		b.addTagOptions,
		b.addLabelOptions,
	}

	args := []string{
		"docker",
		"build",
	}

	args, err := command.NewBuilder(builderFuncs, args).Build()
	if err != nil {
		return nil, err
	}

	return append(args, b.state.Scratch.Workspace()), nil
}

func (b *builder) addDockerfileOptions(cb *command.Builder) error {
	dockerfile, err := b.env.ExpandString(b.task.Dockerfile)
	if err != nil {
		return err
	}

	cb.AddFlagValue("-f", dockerfile)
	return nil
}

func (b *builder) addTagOptions(cb *command.Builder) error {
	for _, tag := range b.task.Tags {
		expanded, err := b.env.ExpandString(tag)
		if err != nil {
			return err
		}

		cb.AddFlagValue("-t", expanded)
	}

	return nil
}

func (b *builder) addLabelOptions(cb *command.Builder) error {
	for _, label := range b.task.Labels {
		expanded, err := b.env.ExpandString(label)
		if err != nil {
			return err
		}

		cb.AddFlagValue("--label", expanded)
	}

	return nil
}
