package runner

import (
	"bytes"
	"io/ioutil"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type loginCommandBuilderState struct {
	state *state.State
	task  *config.LoginTask
	env   environment.Environment
}

func NewLoginCommandRunner(
	state *state.State,
	task *config.LoginTask,
	prefix *logging.Prefix,
	env environment.Environment,
) Runner {
	factory := loginCommandFactory(
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

func loginCommandFactory(
	state *state.State,
	task *config.LoginTask,
	env environment.Environment,
) BuilderFactory {
	return func() (*command.Builder, error) {
		s := &loginCommandBuilderState{
			state: state,
			task:  task,
			env:   env,
		}

		return command.NewBuilder(
			[]string{
				"docker",
				"login",
			},
			[]command.BuildFunc{
				s.addServerArg,
				s.addUsernameOptions,
				s.addPasswordOptions,
			},
		), nil
	}
}

func (s *loginCommandBuilderState) addServerArg(cb *command.Builder) error {
	server, err := s.env.ExpandString(s.task.Server)
	if err != nil {
		return err
	}

	cb.AddArgs(server)
	return nil
}

func (s *loginCommandBuilderState) addUsernameOptions(cb *command.Builder) error {
	username, err := s.env.ExpandString(s.task.Username)
	if err != nil {
		return err
	}

	cb.AddFlagValue("-u", username)
	return nil
}

func (s *loginCommandBuilderState) addPasswordOptions(cb *command.Builder) error {
	password, err := s.env.ExpandString(s.task.Password)
	if err != nil {
		return err
	}

	passwordFile, err := s.env.ExpandString(s.task.PasswordFile)
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

	cb.AddFlag("--password-stdin")
	cb.SetStdin(ioutil.NopCloser(bytes.NewReader([]byte(password))))
	return nil
}
