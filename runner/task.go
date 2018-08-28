package runner

import (
	"github.com/efritz/ij/command"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/state"
)

type (
	Runner interface {
		Run() bool
	}

	baseRunner struct {
		state   *state.State
		prefix  *logging.Prefix
		factory BuilderSetFactory
	}

	BuilderFactory    func() (*command.Builder, error)
	BuilderSetFactory func() ([]*command.Builder, error)
)

func NewBaseRunner(
	state *state.State,
	prefix *logging.Prefix,
	factory BuilderSetFactory,
) Runner {
	return &baseRunner{
		state:   state,
		prefix:  prefix,
		factory: factory,
	}
}

func (r *baseRunner) Run() bool {
	r.state.Logger.Info(
		r.prefix,
		"Beginning task",
	)

	builders, err := r.factory()
	if err != nil {
		r.state.Logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	for _, builder := range builders {
		args, stdin, err := builder.Build()
		if err != nil {
			r.state.Logger.Error(
				r.prefix,
				"Failed to build command args: %s",
				err.Error(),
			)

			return false
		}

		err = command.NewRunner(r.state.Logger).Run(
			r.state.Context,
			args,
			stdin,
			r.prefix,
		)

		if err != nil {
			r.state.ReportError(
				r.prefix,
				"Command failed: %s",
				err.Error(),
			)

			return false
		}
	}

	return true
}

func NewMultiFactory(factory BuilderFactory) BuilderSetFactory {
	return func() ([]*command.Builder, error) {
		builder, err := factory()
		if err != nil {
			return nil, err
		}

		return []*command.Builder{builder}, nil
	}
}
