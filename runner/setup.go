package runner

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/environment"
	"github.com/ij-build/ij/logging"
	"github.com/ij-build/ij/network"
	"github.com/ij-build/ij/options"
	"github.com/ij-build/ij/registry"
	"github.com/ij-build/ij/scratch"
	"github.com/ij-build/ij/ssh"
	"github.com/ij-build/ij/util"
)

func SetupRunner(
	cfg *config.Config,
	appOptions *options.AppOptions,
	runOptions *options.RunOptions,
) (runner *Runner, err error) {
	var (
		cleanup           = NewCleanup()
		ctx, cancel       = setupContext(runOptions.Context, runOptions.PlanTimeout)
		logger            logging.Logger
		loggerFactory     *logging.LoggerFactory
		runID             string
		scratch           *scratch.ScratchSpace
		taskRunnerFactory TaskRunnerFactory
	)

	if runID, err = setupRunID(); err != nil {
		return
	}

	enableHostSSHAgent, err := shouldEnableHostSSHAgent(
		runOptions.EnableContainerSSHAgent,
		cfg.Options.SSHIdentities,
	)

	if err != nil {
		return
	}

	scratch, err = setupScratch(
		runID,
		appOptions.ProjectDir,
		appOptions.ScratchRoot,
		cleanup,
		runOptions.KeepWorkspace,
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

	logger, loggerFactory, err = setupLogger(
		cleanup,
		scratch,
		appOptions.Quiet,
		appOptions.Verbose,
		!appOptions.DisableColor,
		appOptions.FileFactory,
	)

	if err != nil {
		return
	}

	cleanup.Register(func() {
		scratch.Prune(logger)
	})

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

	if runOptions.EnableContainerSSHAgent {
		logger.Info(
			nil,
			"Starting ssh-agent container",
		)

		err = startSSHAgent(
			runID,
			cfg.Options.SSHIdentities,
			scratch,
			containerLists,
			logger,
		)

		if err != nil {
			reportError(
				ctx,
				logger,
				nil,
				"error: failed to setup ssh-agent: %s",
				err.Error(),
			)

			return
		}
	}

	err = setupRegistries(
		ctx,
		cfg,
		cleanup,
		logger,
		appOptions.Env,
		runOptions.Login,
	)

	if err != nil {
		return
	}

	taskRunnerFactory = func(
		context *RunContext,
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
				context,
				t,
				env,
				prefix,
			)

		case *config.RemoveTask:
			return NewRemoveTaskRunnerFactory(
				ctx,
				logger,
			)(
				context,
				t,
				env,
				prefix,
			)

		case *config.RunTask:
			containerOptions := &containerOptions{
				EnableHostSSHAgent:      enableHostSSHAgent,
				EnableContainerSSHAgent: runOptions.EnableContainerSSHAgent,
				CPUShares:               runOptions.CPUShares,
				Memory:                  runOptions.Memory,
			}

			return NewRunTaskRunnerFactory(
				ctx,
				cfg,
				runID,
				scratch,
				containerLists,
				containerOptions,
				logger,
				loggerFactory,
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
				appOptions.Env,
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
		appOptions.Env,
	)

	return
}

//
// Setup Functions

func setupContext(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if timeout == 0 {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, timeout)
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
	projectDir string,
	scratchRoot string,
	cleanup *Cleanup,
	keepWorkspace bool,
) (*scratch.ScratchSpace, error) {
	scratch := scratch.NewScratchSpace(
		runID,
		projectDir,
		scratchRoot,
		keepWorkspace,
	)

	if err := scratch.Setup(); err != nil {
		logging.EmergencyLog(
			"error: failed to create scratch directory: %s",
			err.Error(),
		)

		return nil, err
	}

	return scratch, nil
}

func setupLogger(
	cleanup *Cleanup,
	scratch *scratch.ScratchSpace,
	quiet bool,
	verbose bool,
	colorize bool,
	fileFactory logging.FileFactory,
) (logging.Logger, *logging.LoggerFactory, error) {
	logProcessor := logging.NewProcessor(quiet, verbose, colorize)
	logProcessor.Start()
	cleanup.Register(logProcessor.Shutdown)

	//
	// Create Logger Factory

	if fileFactory == nil {
		fileFactory = func(prefix string) (io.WriteCloser, io.WriteCloser, error) {
			return scratch.MakeLogFiles(prefix)
		}
	}

	loggerFactory := logging.NewLoggerFactory(logProcessor, fileFactory)

	//
	// Create Base Logger

	logger, err := loggerFactory.Logger("ij", true)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to create log files: %s",
			err.Error(),
		)

		return nil, nil, err
	}

	return logger, loggerFactory, nil
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
		reportError(
			ctx,
			logger,
			nil,
			"error: failed to resolve registries: %s",
			err.Error(),
		)

		return err
	}

	if err = registrySet.Login(); err != nil {
		reportError(
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
		reportError(
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

//
// Helpers

func shouldEnableHostSSHAgent(enableContainerSSHAgent bool, identities []string) (bool, error) {
	if enableContainerSSHAgent {
		return false, nil
	}

	enable, err := ssh.EnsureKeysAvailable(identities)
	if err != nil {
		return false, fmt.Errorf(
			"failed to validate ssh keys: %s",
			err.Error(),
		)
	}

	return enable, nil
}

func reportError(
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
