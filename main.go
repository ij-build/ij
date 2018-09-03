package main

import (
	"context"
	"os"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/runner"
	"github.com/efritz/ij/ssh"
	"github.com/efritz/ij/state"
)

const Version = "0.1.0"

var commandRunners = map[string]func(*config.Config) bool{
	"run":    runRun,
	"login":  runLogin,
	"logout": runLogout,
}

func main() {
	if !runMain() {
		os.Exit(1)
	}
}

func runMain() bool {
	command, err := parseArgs()
	if err != nil {
		logging.EmergencyLog("error: %s", err.Error())
		return false
	}

	config, ok := loadConfig()
	if !ok {
		return false
	}

	if !ensureDocker() {
		logging.EmergencyLog("error: docker is not running")
		return false
	}

	if f, ok := commandRunners[command]; ok {
		return f(config)
	}

	panic("unexpected command type")
}

func runRun(cfg *config.Config) bool {
	enableSSHAgent, err := ssh.EnsureKeysAvailable(
		cfg.Options.SSHIdentities,
	)

	if err != nil {
		logging.EmergencyLog(
			"error: failed to validate ssh keys: %s",
			err.Error(),
		)

		return false
	}

	state, err := state.NewState(
		cfg,
		*colorize,
		*cpuShares,
		enableSSHAgent,
		*env,
		*keepWorkspace,
		*loginForPlan,
		*memory,
		*planTimeout,
		*verbose,
	)

	if err != nil {
		return false
	}

	return runner.NewRunner(state, *plans).Run()
}

func runLogin(config *config.Config) bool {
	return withRegistrySet(config, func(registrySet *registry.RegistrySet, logger logging.Logger) bool {
		if err := registrySet.Login(); err != nil {
			logger.Error(nil, "failed to log in to registries: %s", err.Error())
			return false
		}

		return true
	})
}

func runLogout(config *config.Config) bool {
	return withRegistrySet(config, func(registrySet *registry.RegistrySet, logger logging.Logger) bool {
		registrySet.Logout()
		return true
	})
}

func withRegistrySet(config *config.Config, f func(*registry.RegistrySet, logging.Logger) bool) bool {
	logProcessor := logging.NewProcessor(*verbose, *colorize)
	logProcessor.Start()
	defer logProcessor.Shutdown()

	logger := logProcessor.Logger(
		logging.NilWriter,
		logging.NilWriter,
		true,
	)

	registryEnv := environment.Merge(
		environment.New(config.Environment),
		environment.New(*env),
	)

	registrySet, err := registry.NewRegistrySet(
		context.Background(),
		logger,
		registryEnv,
		config.Registries,
	)

	if err != nil {
		logger.Error(nil, "failed to create registry set: %s", err.Error())
		return false
	}

	return f(registrySet, logger)
}
