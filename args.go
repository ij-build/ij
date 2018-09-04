package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/efritz/ij/paths"
)

var (
	app = kingpin.New("ij", "IJ is a build tool using Docker containers.").Version(Version)

	// Commands
	login      = app.Command("login", "Login to docker registries.")
	logout     = app.Command("logout", "Logout of docker registries.")
	rotateLogs = app.Command("rotate-logs", "Trim old run logs the .ij directory.")
	run        = app.Command("run", "Run a plan or metaplan.").Default()

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
)

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
