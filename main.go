package main

import (
	"context"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
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
		"ij-override.yaml",
		"ij-override.yml",
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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)

	defer cancel()

	if !ensureDocker(ctx) {
		logging.EmergencyLog("error: docker is not running")
		return false
	}

	switch command {
	case "run":
		return runRun(ctx, config)
	case "login":
		return runLogin(ctx, config)
	case "logout":
		return runLogout(ctx, config)
	}

	panic("unexpected command type")
}

func runRun(ctx context.Context, config *config.Config) bool {
	enableSSHAgent, err := ssh.EnsureKeysAvailable(append(
		config.SSHIdentities,
		*sshIdentities...,
	))

	if err != nil {
		logging.EmergencyLog(
			"error: failed to validate ssh keys: %s",
			err.Error(),
		)

		return false
	}

	state, err := state.NewState(
		config,
		*plans,
		*colorize,
		*cpuShares,
		enableSSHAgent,
		*env,
		*forceSequential,
		*healthcheckInterval,
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

func runLogin(ctx context.Context, config *config.Config) bool {
	// TODO - implement these commands
	return false
}

func runLogout(ctx context.Context, config *config.Config) bool {
	// TODO - implement these commands
	return false
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
	}

	return command, nil
}

func loadConfig() (*config.Config, bool) {
	loader := loader.NewLoader()

	config, err := loader.Load(*configPath)
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

	if err := loader.ApplyOverrides(config, overridePaths); err != nil {
		logging.EmergencyLog(
			"error: failed to apply overrides: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := config.Resolve(); err != nil {
		logging.EmergencyLog(
			"error: failed to resolve config: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := config.Validate(); err != nil {
		logging.EmergencyLog(
			"error: failed to validate config: %s",
			err.Error(),
		)

		return nil, false
	}

	for _, name := range *plans {
		if !config.IsPlanDefined(name) {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return nil, false
		}
	}

	return config, true
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

func ensureDocker(ctx context.Context) bool {
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
