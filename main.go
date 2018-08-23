package main

import (
	"context"
	"os"
	"time"

	"github.com/alecthomas/kingpin"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
	"github.com/efritz/ij/runtime"
	"github.com/efritz/ij/ssh"
)

const Version = "0.1.0"

var (
	app = kingpin.New("ij", "ij is a build tool built with Docker.").Version(Version)

	plans               = app.Arg("plans", "The name of the plans to execute.").Default("default").Strings()
	configPath          = app.Flag("config", "The path to the config file.").Short('f').String()
	env                 = app.Flag("env", "Environment variables.").Short('e').Strings()
	verbose             = app.Flag("verbose", "Output debug logs.").Short('v').Default("false").Bool()
	colorize            = app.Flag("color", "Enable colorized output.").Default("true").Bool()
	forceSequential     = app.Flag("force-sequential", "Disable parallel execution.").Default("false").Bool()
	healthcheckInterval = app.Flag("healthcheck-interval", "The interval between service container healthchecks.").Default("5s").Duration()
	cpuShares           = app.Flag("cpu-shares", "The amount of cpu shares to give to each container.").Short('c').String()
	memory              = app.Flag("memory", "The amount of memory to give each container.").Short('m').String()
	sshIdentities       = app.Flag("ssh-identity", "Enable ssh-agent for the given identities.").Strings()
	planTimeout         = app.Flag("timeout", "Maximum amount of time a plan can run. 0 to disable.").Default("15m").Duration()

	defaultConfigPaths = []string{
		"ij.yaml",
		"ij.yml",
	}
)

func main() {
	if !run() {
		os.Exit(1)
	}
}

func run() bool {
	if err := parseArgs(); err != nil {
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

	enableAgent, err := ssh.EnureKeysAvailable(*sshIdentities)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to validate ssh keys: %s",
			err.Error(),
		)

		return false
	}

	state, err := runtime.NewState(
		config,
		*plans,
		*env,
		*verbose,
		*colorize,
		*forceSequential,
		enableAgent,
		*healthcheckInterval,
		*cpuShares,
		*memory,
		*planTimeout,
	)

	if err != nil {
		return false
	}

	return runtime.NewPlanRunner(state).Run()
}

func parseArgs() error {
	args := os.Args[1:]

	if _, err := app.Parse(args); err != nil {
		return err
	}

	if *configPath == "" {
		for _, path := range defaultConfigPaths {
			ok, err := paths.FileExists(path)
			if err != nil {
				return err
			}

			if ok {
				*configPath = path
				break
			}
		}
	}

	return nil
}

func loadConfig() (*config.Config, bool) {
	config, err := loader.NewLoader().Load(*configPath)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to load config: %s",
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
		if _, ok := config.Plans[name]; !ok {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return nil, false
		}
	}

	return config, true
}

func ensureDocker(ctx context.Context) bool {
	_, _, err := command.RunForOutput(
		ctx,
		[]string{
			"docker",
			"ps",
			"-q",
		},
		nil,
	)

	return err == nil
}
