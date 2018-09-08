package runner

import (
	"context"
	"time"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/network"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/scratch"
	"github.com/efritz/ij/util"
)

func SetupRunner(
	cfg *config.Config,
	colorize bool,
	overrideEnv []string,
	verbose bool,
	enableSSHAgent bool,
	cpuShares string,
	keepWorkspace bool,
	login bool,
	memory string,
	planTimeout time.Duration,
) (runner *Runner, err error) {
	var (
		cleanup           = NewCleanup()
		ctx, cancel       = setupContext(planTimeout)
		logger            logging.Logger
		logProcessor      logging.Processor
		runID             string
		scratch           *scratch.ScratchSpace
		taskRunnerFactory TaskRunnerFactory
	)

	if runID, err = setupRunID(); err != nil {
		return
	}

	scratch, err = setupScratch(
		runID,
		cleanup,
		keepWorkspace,
	)

	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			return
		}

		if err := scratch.Teardown(); err != nil {
			logging.EmergencyLog(
				"error: failed to teardown scratch directory: %s",
				err.Error(),
			)
		}
	}()

	logProcessor, logger, err = setupLogger(
		cleanup,
		scratch,
		verbose,
		colorize,
	)

	if err != nil {
		return
	}

	_, err = setupNetwork(
		ctx,
		runID,
		cleanup,
		logger,
	)

	if err != nil {
		return
	}

	containerLists := setupContainerLists(
		runID,
		cleanup,
		logger,
	)

	err = setupRegistries(
		ctx,
		cfg,
		cleanup,
		logger,
		overrideEnv,
		login,
	)

	if err != nil {
		return
	}

	taskRunnerFactory = func(
		task config.Task,
		prefix *logging.Prefix,
		env environment.Environment,
	) TaskRunner {
		switch t := task.(type) {
		case *config.BuildTask:
			return NewBuildTaskRunnerFactory(
				ctx,
				scratch.Workspace(),
				logger,
			)(
				t,
				env,
				prefix,
			)

		case *config.PushTask:
			return NewPushTaskRunnerFactory(
				ctx,
				logger,
			)(
				t,
				env,
				prefix,
			)

		case *config.RemoveTask:
			return NewRemoveTaskRunnerFactory(
				ctx,
				logger,
			)(
				t,
				env,
				prefix,
			)

		case *config.RunTask:
			containerOptions := &containerOptions{
				EnableSSHAgent: enableSSHAgent,
				CPUShares:      cpuShares,
				Memory:         memory,
			}

			return NewRunTaskRunnerFactory(
				ctx,
				cfg,
				runID,
				scratch,
				containerLists,
				containerOptions,
				logProcessor,
				logger,
			)(
				t,
				env,
				prefix,
			)

		case *config.PlanTask:
			runner := NewPlanRunner(
				ctx,
				cfg,
				taskRunnerFactory,
				logger,
				overrideEnv,
			)

			return NewPlanTaskRunnerFactory(
				runner,
				logger,
			)(
				t,
				env,
				prefix,
			)
		}

		panic("unexpected task type")
	}

	runner = NewRunner(
		ctx,
		logger,
		cfg,
		taskRunnerFactory,
		scratch,
		cleanup,
		runID,
		cancel,
		overrideEnv,
	)

	return
}

//
// Setup Functions

func setupContext(timeout time.Duration) (context.Context, func()) {
	if timeout == 0 {
		return context.WithCancel(context.Background())
	}

	return context.WithTimeout(context.Background(), timeout)
}

func setupRunID() (string, error) {
	id, err := util.MakeID()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to generate run id: %s",
			err.Error(),
		)

		return "", err
	}

	return id, nil
}

func setupScratch(
	runID string,
	cleanup *Cleanup,
	keepWorkspace bool,
) (*scratch.ScratchSpace, error) {
	scratch := scratch.NewScratchSpace(runID, keepWorkspace)

	if err := scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create scratch directory: %s",
			err.Error(),
		)

		return nil, err
	}

	cleanup.Register(func() {
		if err := scratch.Prune(); err != nil {
			logging.EmergencyLog(
				"error: failed to clean up scratch directory: %s",
				err.Error(),
			)
		}
	})

	return scratch, nil
}

func setupLogger(
	cleanup *Cleanup,
	scratch *scratch.ScratchSpace,
	verbose bool,
	colorize bool,
) (logging.Processor, logging.Logger, error) {
	logProcessor := logging.NewProcessor(verbose, colorize)
	logProcessor.Start()
	cleanup.Register(logProcessor.Shutdown)

	//
	// Create Base Logger

	outfile, errfile, err := scratch.MakeLogFiles("ij")
	if err != nil {
		logging.EmergencyLog(
			"error: failed to create log files: %s",
			err.Error(),
		)

		return nil, nil, err
	}

	logger := logProcessor.Logger(
		outfile,
		errfile,
		true,
	)

	return logProcessor, logger, nil
}

func setupContainerLists(
	runID string,
	cleanup *Cleanup,
	logger logging.Logger,
) *ContainerLists {
	containerStopper := NewContainerStopper(logger)
	networkDisconnector := NewNetworkDisconnector(runID, logger)

	cleanup.Register(containerStopper.Execute)
	cleanup.Register(networkDisconnector.Execute)

	return &ContainerLists{
		ContainerStopper:    containerStopper,
		NetworkDisconnector: networkDisconnector,
	}
}

func setupRegistries(
	ctx context.Context,
	cfg *config.Config,
	cleanup *Cleanup,
	logger logging.Logger,
	env []string,
	login bool,
) error {
	if !login {
		return nil
	}

	registryEnv := environment.Merge(
		environment.New(cfg.Environment),
		environment.New(env),
	)

	registrySet, err := registry.NewRegistrySet(
		ctx,
		logger,
		registryEnv,
		cfg.Registries,
	)

	if err != nil {
		ReportError(
			ctx,
			logger,
			nil,
			"error: failed to resolve registries: %s",
			err.Error(),
		)

		return err
	}

	if err = registrySet.Login(); err != nil {
		ReportError(
			ctx,
			logger,
			nil,
			"error: failed to log into registries: %s",
			err.Error(),
		)

		return err
	}

	cleanup.Register(registrySet.Logout)
	return nil
}

func setupNetwork(
	ctx context.Context,
	runID string,
	cleanup *Cleanup,
	logger logging.Logger,
) (*network.Network, error) {
	network, err := network.NewNetwork(ctx, runID, logger)
	if err != nil {
		ReportError(
			ctx,
			logger,
			nil,
			"error: failed to create network: %s",
			err.Error(),
		)

		return nil, err
	}

	cleanup.Register(network.Teardown)
	return network, nil
}

func ReportError(
	ctx context.Context,
	logger logging.Logger,
	prefix *logging.Prefix,
	format string,
	args ...interface{},
) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	logger.Error(prefix, format, args...)
}
