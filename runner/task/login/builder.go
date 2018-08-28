package login

import (
	"io/ioutil"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/state"
)

type builder struct {
	state *state.State
	task  *config.LoginTask
	env   environment.Environment
}

func build(
	state *state.State,
	task *config.LoginTask,
	env environment.Environment,
) ([]string, error) {
	b := &builder{
		state: state,
		task:  task,
		env:   env,
	}

	builderFuncs := []command.BuildFunc{
		b.addUsernameOptions,
		b.addPasswordOptions,
		// TODO - is last
		b.addServerOptions,
	}

	args := []string{
		"docker",
		"login",
	}

	args, err := command.NewBuilder(builderFuncs, args).Build()
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (b *builder) addServerOptions(cb *command.Builder) error {
	server, err := b.env.ExpandString(b.task.Server)
	if err != nil {
		return err
	}

	cb.AddFlag(server)
	return nil
}

func (b *builder) addUsernameOptions(cb *command.Builder) error {
	username, err := b.env.ExpandString(b.task.Username)
	if err != nil {
		return err
	}

	cb.AddFlagValue("-u", username)
	return nil
}

func (b *builder) addPasswordOptions(cb *command.Builder) error {
	password, err := b.env.ExpandString(b.task.Password)
	if err != nil {
		return err
	}

	passwordFile, err := b.env.ExpandString(b.task.PasswordFile)
	if err != nil {
		return err
	}

	if passwordFile != "" {
		content, err := ioutil.ReadFile(passwordFile)
		if err != nil {
			return err
		}

		password = string(content)
	}

	// TODO - use --password-stdin instead
	cb.AddFlagValue("--password", password)
	return nil
}
