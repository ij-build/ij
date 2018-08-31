package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/loader"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/registry"
	"github.com/efritz/ij/runner"
	"github.com/efritz/ij/ssh"
	"github.com/efritz/ij/state"
)

const Version = "0.1.0"

var (
	app = kingpin.New("ij", "ij is a build tool built with Docker.").Version(Version)

	// Commands
	run    = app.Command("run", "Run a plan or metaplan.").Default()
	login  = app.Command("login", "Login to docker registries.")
	logout = app.Command("logout", "Logout of docker registries.")

	// Shared options
	colorize   = app.Flag("color", "Enable colorized output.").Default("true").Bool()
	configPath = app.Flag("config", "The path to the config file.").Short('f').String()
	env        = app.Flag("env", "Environment variables.").Short('e').Strings()
	envFiles   = app.Flag("env-file", "Environment file.").Strings()
	verbose    = app.Flag("verbose", "Output debug logs.").Short('v').Default("false").Bool()

	// Run Options
	plans               = run.Arg("plans", "The name of the plans to execute.").Default("default").Strings()
	cpuShares           = run.Flag("cpu-shares", "The amount of cpu shares to give to each container.").Short('c').String()
	forceSequential     = run.Flag("force-sequential", "Disable parallel execution.").Default("false").Bool()
	healthcheckInterval = run.Flag("healthcheck-interval", "The interval between service container healthchecks.").Default("5s").Duration()
	keepWorkspace       = run.Flag("keep-workspace", "Do not delete the workspace").Short('k').Default("false").Bool()
	loginForPlan        = run.Flag("login", "Login to docker registries before running.").Default("false").Bool()
	memory              = run.Flag("memory", "The amount of memory to give each container.").Short('m').String()
	planTimeout         = run.Flag("timeout", "Maximum amount of time a plan can run. 0 to disable.").Default("15m").Duration()
	sshIdentities       = run.Flag("ssh-identity", "Enable ssh-agent for the given identities.").Strings()

	defaultConfigPaths = []string{
		"ij.yaml",
		"ij.yml",
	}

	localOverridePaths = []string{
		"ij.override.yaml",
		"ij.override.yml",
	}

	globalOverridePaths = []string{
		filepath.Join(".ij", "override.yaml"),
	}
)

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

	switch command {
	case "run":
		return runRun(config)
	case "login":
		return runLogin(config)
	case "logout":
		return runLogout(config)
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
		*plans,
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

	return runner.NewPlanRunner(state).Run()
}

func runLogin(config *config.Config) bool {
	return withRegistrySet(config, func(registrySet *registry.RegistrySet, logger logging.Logger) bool {
		if err := registrySet.Login(); err != nil {
			logger.Error(nil, "failed to login to registries: %s", err.Error())
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
		environment.Default(),
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

func parseArgs() (string, error) {
	command, err := app.Parse(os.Args[1:])
	if err != nil {
		return "", err
	}

	if *configPath == "" {
		for _, path := range defaultConfigPaths {
			ok, err := paths.FileExists(path)
			if err != nil {
				return "", err
			}

			if ok {
				*configPath = path
				break
			}
		}

		if *configPath == "" {
			return "", fmt.Errorf("could not infer config file")
		}
	}

	return command, nil
}

func loadConfig() (*config.Config, bool) {
	loader := loader.NewLoader()

	cfg, err := loader.Load(*configPath)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to load config: %s",
			err.Error(),
		)

		return nil, false
	}

	overridePaths, err := getOverridePaths()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to determine override paths: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := loader.ApplyOverrides(cfg, overridePaths); err != nil {
		logging.EmergencyLog(
			"error: failed to apply overrides: %s",
			err.Error(),
		)

		return nil, false
	}

	cfg.ApplyOverride(&config.Override{
		Options: &config.Options{
			SSHIdentities:       *sshIdentities,
			ForceSequential:     *forceSequential,
			HealthcheckInterval: *healthcheckInterval,
		},
		EnvironmentFiles: *envFiles,
	})

	envFromFile := []string{}
	for _, path := range cfg.EnvironmentFiles {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			logging.EmergencyLog(
				"error: failed to read environment file: %s",
				err.Error(),
			)

			return nil, false
		}

		lines, err := environment.NormalizeEnvironmentFile(string(content))
		if err != nil {
			logging.EmergencyLog(
				"error: failed to read environment file: %s",
				err.Error(),
			)

			return nil, false
		}

		envFromFile = append(envFromFile, lines...)
	}

	*env = append(envFromFile, *env...)

	if err := cfg.Resolve(); err != nil {
		logging.EmergencyLog(
			"error: failed to resolve config: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := cfg.Validate(); err != nil {
		logging.EmergencyLog(
			"error: failed to validate config: %s",
			err.Error(),
		)

		return nil, false
	}

	for _, name := range *plans {
		if !cfg.IsPlanDefined(name) {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return nil, false
		}
	}

	return cfg, true
}

func getOverridePaths() ([]string, error) {
	current, err := user.Current()
	if err != nil {
		return nil, err
	}

	found := []string{}
	for _, path := range globalOverridePaths {
		path = filepath.Join(current.HomeDir, path)

		if ok, err := paths.FileExists(path); err == nil && ok {
			found = append(found, path)
			break
		}
	}

	for _, path := range localOverridePaths {
		if ok, err := paths.FileExists(path); err == nil && ok {
			found = append(found, path)
		}
	}

	return found, nil
}

func ensureDocker() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := []string{
		"docker",
		"ps",
		"-q",
	}

	_, _, err := command.NewRunner(nil).RunForOutput(
		ctx,
		args,
		nil,
	)

	return err == nil
}
