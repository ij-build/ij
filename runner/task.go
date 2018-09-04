package runner

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/logging"
)

type (
	TaskRunner interface {
		Run(*RunContext) bool
	}

	baseRunner struct {
		ctx     context.Context
		factory BuilderSetFactory
		logger  logging.Logger
		prefix  *logging.Prefix
	}

	BuilderFactory    func() (*command.Builder, error)
	BuilderSetFactory func() ([]*command.Builder, error)
)

func NewBaseRunner(
	ctx context.Context,
	factory BuilderSetFactory,
	logger logging.Logger,
	prefix *logging.Prefix,
) TaskRunner {
	return &baseRunner{
		ctx:     ctx,
		logger:  logger,
		prefix:  prefix,
		factory: factory,
	}
}

func (r *baseRunner) Run(context *RunContext) bool {
	r.logger.Info(
		r.prefix,
		"Beginning task",
	)

	builders, err := r.factory()
	if err != nil {
		r.logger.Error(
			r.prefix,
			"Failed to build command args: %s",
			err.Error(),
		)

		return false
	}

	for _, builder := range builders {
		args, stdin, err := builder.Build()
		if err != nil {
			r.logger.Error(
				r.prefix,
				"Failed to build command args: %s",
				err.Error(),
			)

			return false
		}

		err = command.NewRunner(r.logger).Run(
			r.ctx,
			args,
			stdin,
			r.prefix,
		)

		if err != nil {
			ReportError(
				r.ctx,
				r.logger,
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
