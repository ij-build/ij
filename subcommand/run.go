package subcommand

import (
	"context"
	"fmt"
	"time"

	"github.com/efritz/ij/command"
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/runner"
	"github.com/efritz/ij/ssh"
)

type RunOptions struct {
	Plans               []string
	CPUShares           string
	ForceSequential     bool
	HealthcheckInterval time.Duration
	KeepWorkspace       bool
	LoginForPlan        bool
	Memory              string
	PlanTimeout         time.Duration
	SSHIdentities       []string
}

var ErrFailed = fmt.Errorf("subcommand failed")

func NewRunCommand(appOptions *AppOptions, runOptions *RunOptions) CommandRunner {
	return func(config *config.Config) error {
		if !ensureDocker() {
			return fmt.Errorf("docker is not running")
		}

		enableSSHAgent, err := ssh.EnsureKeysAvailable(
			config.Options.SSHIdentities,
		)

		if err != nil {
			return fmt.Errorf(
				"failed to validate ssh keys: %s",
				err.Error(),
			)
		}

		for _, name := range runOptions.Plans {
			if !config.IsPlanDefined(name) {
				return fmt.Errorf(
					"unknown plan %s",
					name,
				)
			}
		}

		runner, err := runner.SetupRunner(
			config,
			appOptions.Colorize,
			appOptions.Env,
			appOptions.Verbose,
			enableSSHAgent,
			runOptions.CPUShares,
			runOptions.KeepWorkspace,
			runOptions.LoginForPlan,
			runOptions.Memory,
			runOptions.PlanTimeout,
		)

		if err != nil {
			return err
		}

		if !runner.Run(runOptions.Plans) {
			return ErrFailed
		}

		return nil
	}
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
