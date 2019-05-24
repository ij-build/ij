package runner

import (
	"context"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
)

type (
	TaskRunner interface {
		Run(*RunContext) bool
	}

	BaseRunner interface {
		TaskRunner
		RegisterOnSuccess(HookFunc)
		RegisterOnFailure(HookFunc)
	}

	TaskRunnerFactory func(
		*RunContext,
		config.Task,
		*logging.Prefix,
		environment.Environment,
	) TaskRunner

	baseRunner struct {
		ctx       context.Context
		factory   BuilderSetFactory
		logger    logging.Logger
		prefix    *logging.Prefix
		onSuccess HookFunc
		onFailure HookFunc
	}

	HookFunc          func(context *RunContext) error
	BuilderFactory    func() (*command.Builder, error)
	BuilderSetFactory func() ([]*command.Builder, error)
)

func NewBaseRunner(
	ctx context.Context,
	factory BuilderSetFactory,
	logger logging.Logger,
	prefix *logging.Prefix,
) BaseRunner {
	return &baseRunner{
		ctx:       ctx,
		logger:    logger,
		prefix:    prefix,
		factory:   factory,
		onSuccess: func(context *RunContext) error { return nil },
		onFailure: func(context *RunContext) error { return nil },
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

		return r.runFailureHook(context)
	}

	for _, builder := range builders {
		args, stdin, err := builder.Build()
		if err != nil {
			r.logger.Error(
				r.prefix,
				"Failed to build command args: %s",
				err.Error(),
			)

			return r.runFailureHook(context)
		}

		err = command.NewRunner(r.logger).Run(
			r.ctx,
			args,
			stdin,
			r.prefix,
		)

		if err != nil {
			reportError(
				r.ctx,
				r.logger,
				r.prefix,
				"Command failed: %s",
				err.Error(),
			)

			return r.runFailureHook(context)
		}
	}

	return r.runSuccessHook(context)
}

func (r *baseRunner) RegisterOnSuccess(hookFunc HookFunc) {
	r.onSuccess = hookFunc
}

func (r *baseRunner) RegisterOnFailure(hookFunc HookFunc) {
	r.onFailure = hookFunc
}

func (r *baseRunner) runSuccessHook(context *RunContext) bool {
	if err := r.onSuccess(context); err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Success hook failed: %s",
			err.Error(),
		)

		return false
	}

	return true
}

func (r *baseRunner) runFailureHook(context *RunContext) bool {
	if err := r.onFailure(context); err != nil {
		reportError(
			r.ctx,
			r.logger,
			r.prefix,
			"Failure hook failed: %s",
			err.Error(),
		)
	}

	return false
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
