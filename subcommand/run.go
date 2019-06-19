package subcommand

import (
	"context"
	"fmt"
	"time"

	"github.com/ij-build/ij/command"
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/options"
	"github.com/ij-build/ij/runner"
)

var ErrBuildFailed = fmt.Errorf("subcommand failed")

func NewRunCommand(appOptions *options.AppOptions, runOptions *options.RunOptions) CommandRunner {
	return func(config *config.Config) error {
		if !ensureDocker() {
			return fmt.Errorf("docker is not running")
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
			appOptions,
			runOptions,
		)

		if err != nil {
			return err
		}

		if !runner.Run(runOptions.Plans) {
			return ErrBuildFailed
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
