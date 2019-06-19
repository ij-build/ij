package subcommand

import (
	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/options"
)

type CommandRunner func(*config.Config) error

func Run(
	command string,
	config *config.Config,
	appOptions *options.AppOptions,
	cleanOptions *options.CleanOptions,
	runOptions *options.RunOptions,
) error {
	runners := map[string]CommandRunner{
		"clean":       NewCleanCommand(appOptions, cleanOptions),
		"login":       NewLoginCommand(appOptions),
		"logout":      NewLogoutCommand(appOptions),
		"rotate-logs": NewRotateLogsCommand(appOptions),
		"run":         NewRunCommand(appOptions, runOptions),
		"show-config": NewShowConfigCommand(appOptions),
	}

	runner, ok := runners[command]
	if !ok {
		panic("unexpected command type")
	}

	return runner(config)
}
