package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/consts"
	"github.com/ij-build/ij/loader"
	"github.com/ij-build/ij/logging"
	"github.com/ij-build/ij/options"
	"github.com/ij-build/ij/subcommand"
)

func newSharedOptions(app *kingpin.Application, projectDir string) *options.AppOptions {
	opts := &options.AppOptions{
		ProjectDir:  projectDir,
		ScratchRoot: projectDir,
	}

	app.Flag("scratch-root", "The directory where .ij/ files are written.").StringVar(&opts.ScratchRoot)
	app.Flag("config", "The path to the config file.").Short('f').StringVar(&opts.ConfigPath)
	app.Flag("env", "Environment variables.").Short('e').StringsVar(&opts.Env)
	app.Flag("env-file", "Environment file.").StringsVar(&opts.EnvFiles)
	app.Flag("quiet", "Do not output to stdout or stderr.").Short('q').Default("false").BoolVar(&opts.Quiet)
	app.Flag("verbose", "Output debug logs.").Short('v').Default("false").BoolVar(&opts.Verbose)
	app.Flag("no-color", "Disable colorized output.").Default("false").BoolVar(&opts.DisableColor)
	return opts
}

func newRunOptions(cmd *kingpin.CmdClause) *options.RunOptions {
	opts := &options.RunOptions{
		Context: context.Background(),
	}

	cmd.Arg("plans", "The name of the plans to execute.").Default("default").StringsVar(&opts.Plans)
	cmd.Flag("cpu-shares", "The amount of cpu shares to give to each container.").Short('c').StringVar(&opts.CPUShares)
	cmd.Flag("force-sequential", "Disable parallel execution.").Default("false").BoolVar(&opts.ForceSequential)
	cmd.Flag("healthcheck-interval", "The interval between service container healthchecks.").Default("5s").DurationVar(&opts.HealthcheckInterval)
	cmd.Flag("keep-workspace", "Do not delete the workspace").Short('k').Default("false").BoolVar(&opts.KeepWorkspace)
	cmd.Flag("login", "Login to docker registries before running.").Default("false").BoolVar(&opts.Login)
	cmd.Flag("memory", "The amount of memory to give each container.").Short('m').StringVar(&opts.Memory)
	cmd.Flag("timeout", "Maximum amount of time a plan can run. 0 to disable.").Default("15m").DurationVar(&opts.PlanTimeout)
	cmd.Flag("ssh-identity", "Enable ssh-agent for the given identities.").StringsVar(&opts.SSHIdentities)
	cmd.Flag("ssh-agent-container", "Start an ssh-agent inside of a container.").BoolVar(&opts.EnableContainerSSHAgent)
	return opts
}

func newCleanOptions(cmd *kingpin.CmdClause) *options.CleanOptions {
	opts := &options.CleanOptions{}
	cmd.Flag("force", "Do not require confirmation before removing matching files.").Default("false").BoolVar(&opts.ForceClean)
	return opts
}

func main() {
	if err := runMain(); err != nil {
		if err != subcommand.ErrBuildFailed {
			logging.EmergencyLog("error: %s", err.Error())
		}

		os.Exit(1)
	}
}

func runMain() error {
	app := kingpin.New("ij", "IJ is a build tool using Docker containers.").Version(consts.Version)
	clean := app.Command("clean", "Remove exported files.")
	_ = app.Command("login", "Login to docker registries.")
	_ = app.Command("logout", "Logout of docker registries.")
	_ = app.Command("rotate-logs", "Trim old run logs the .ij directory.")
	run := app.Command("run", "Run a plan or metaplan.").Default()
	_ = app.Command("show-config", "Show the effective config after resolving parents.")

	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory (%s)", err.Error())
	}

	appOptions := newSharedOptions(app, projectDir)
	cleanOptions := newCleanOptions(clean)
	runOptions := newRunOptions(run)

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	path, err := loader.GetConfigPath(appOptions.ConfigPath)
	if err != nil {
		return err
	}

	override := &config.Override{
		Options: &config.Options{
			SSHIdentities:       runOptions.SSHIdentities,
			ForceSequential:     runOptions.ForceSequential,
			HealthcheckInterval: runOptions.HealthcheckInterval,
		},
		EnvironmentFiles: appOptions.EnvFiles,
	}

	config, err := loader.LoadFile(path, override)
	if err != nil {
		return err
	}

	return subcommand.Run(
		command,
		config,
		appOptions,
		cleanOptions,
		runOptions,
	)
}
