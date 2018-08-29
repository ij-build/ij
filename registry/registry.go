package registry

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/efritz/ij/command"
)

type Login interface {
	GetServer() (string, error)
	Login() error
}

func login(
	ctx context.Context,
	runner command.Runner,
	server string,
	username string,
	password string,
) error {
	builder := command.NewBuilder([]string{
		"docker",
		"login",
	}, nil)

	builder.AddArgs(server)
	builder.AddFlagValue("-u", username)
	builder.AddFlag("--password-stdin")
	builder.SetStdin(ioutil.NopCloser(bytes.NewReader([]byte(password))))

	args, stdin, err := builder.Build()
	if err != nil {
		return err
	}

	return runner.Run(
		ctx,
		args,
		stdin,
		nil,
	)
}
