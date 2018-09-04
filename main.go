package main

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/runner"
	"github.com/efritz/ij/scratch"
	"github.com/efritz/ij/ssh"
)

const Version = "0.1.0"

var commandRunners = map[string]func(*config.Config) bool{
	"login":       runLogin,
	"logout":      runLogout,
	"rotate-logs": runRotateLogs,
	"run":         runRun,
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

func runRotateLogs(config *config.Config) bool {
	wd, err := os.Getwd()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to get working directory: %s",
			err.Error(),
		)

		return false
	}

	scratchPath := filepath.Join(wd, scratch.ScratchDir)

	entries, err := paths.DirContents(scratchPath)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to read scratch directory: %s",
			err.Error(),
		)

		return false
	}

	if len(entries) == 0 {
		return true
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModTime().After(entries[j].ModTime())
	})

	for _, info := range entries[1:] {
		if err := os.RemoveAll(filepath.Join(scratchPath, info.Name())); err != nil {
			logging.EmergencyLog(
				"error: failed to delete run directory: %s",
				err.Error(),
			)

		}
	}

	return true
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

	runner, err := runner.SetupRunner(
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

	return runner.Run(*plans)
}

//
// Context Helpers

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
