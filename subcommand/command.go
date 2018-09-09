package subcommand

import (
	"github.com/efritz/ij/config"
	"github.com/efritz/ij/options"
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
		"clean":       NewCleanCommand(cleanOptions),
		"login":       NewLoginCommand(appOptions),
		"logout":      NewLogoutCommand(appOptions),
		"rotate-logs": NewRotateLogsCommand(),
		"run":         NewRunCommand(appOptions, runOptions),
	}

	runner, ok := runners[command]
	if !ok {
		panic("unexpected command type")
	}

	return runner(config)
}
