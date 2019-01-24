package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/scratch"
)

type Runner struct {
	ctx               context.Context
	logger            logging.Logger
	config            *config.Config
	taskRunnerFactory TaskRunnerFactory
	scratch           *scratch.ScratchSpace
	cleanup           *Cleanup
	runID             string
	cancel            func()
	env               []string
}

const FlashPermissionsImage = "alpine:3.8"

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func NewRunner(
	ctx context.Context,
	logger logging.Logger,
	config *config.Config,
	taskRunnerFactory TaskRunnerFactory,
	scratch *scratch.ScratchSpace,
	cleanup *Cleanup,
	runID string,
	cancel func(),
	env []string,
) *Runner {
	return &Runner{
		ctx:               ctx,
		logger:            logger,
		config:            config,
		taskRunnerFactory: taskRunnerFactory,
		scratch:           scratch,
		cleanup:           cleanup,
		runID:             runID,
		cancel:            cancel,
		env:               env,
	}
}

func (r *Runner) Run(plans []string) bool {
	r.logger.Info(
		nil,
		"Beginning run %s",
		r.runID,
	)

	go r.watchSignals()
	defer r.cleanup.Cleanup()

	defer func() {
		r.logger.Info(
			nil,
			"Finished run %s",
			r.runID,
		)
	}()

	transferer := paths.NewTransferer(
		r.scratch.Project(),
		r.scratch.Scratch(),
		r.scratch.Workspace(),
		r.logger,
	)

	r.logger.Info(
		nil,
		"Importing files to workspace",
	)

	importErr := transferer.Import(
		r.config.Import.Files,
		r.config.Import.Excludes,
	)

	if importErr != nil {
		r.logger.Error(
			nil,
			"Failed to import files to workspace: %s",
			importErr.Error(),
		)

		return false
	}

	var (
		failure     = false
		rootContext = NewRunContext(nil)
	)

	for _, name := range plans {
		runner := NewPlanRunner(
			r.ctx,
			r.config,
			r.taskRunnerFactory,
			r.logger,
			r.env,
		)

		newContext := NewRunContext(rootContext)
		newContext.Failure = failure

		if !runner.Run(newContext, name, logging.NewPrefix()) {
			failure = true
		}
	}

	r.logger.Info(
		nil,
		"Flashing workspace permissions",
	)

	if err := r.flashPermissions(); err != nil {
		r.logger.Error(
			nil,
			"Failed to flash workspace permissions: %s",
			err.Error(),
		)
	}

	if failure {
		return false
	}

	r.logger.Info(
		nil,
		"Exporting files from workspace",
	)

	exportErr := transferer.Export(
		r.config.Export.Files,
		r.config.Export.Excludes,
	)

	if exportErr != nil {
		r.logger.Error(
			nil,
			"Failed to export files from workspace: %s",
			exportErr.Error(),
		)

		return false
	}

	return true
}

func (r *Runner) watchSignals() {
	signals := make(chan os.Signal, 1)

	for _, s := range shutdownSignals {
		signal.Notify(signals, s)
	}

	for range signals {
		r.logger.Error(
			nil,
			"Received signal",
		)

		r.cancel()
		return
	}
}

func (r *Runner) flashPermissions() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	builder := command.NewBuilder([]string{
		"docker",
		"run",
		"--rm",
	}, nil)

	builder.AddArgs(FlashPermissionsImage)
	builder.AddArgs("chown", fmt.Sprintf("%s:%s", user.Uid, user.Gid), "-R", ".")
	builder.AddFlagValue("-w", "/workspace")
	builder.AddFlagValue("-v", fmt.Sprintf("%s:/workspace", r.scratch.Workspace()))

	args, _, err := builder.Build()
	if err != nil {
		return err
	}

	return command.NewRunner(r.logger).Run(
		r.ctx,
		args,
		nil,
		nil,
	)
}
