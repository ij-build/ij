package subcommand

import (
	"github.com/efritz/ij/config"
)

type (
	AppOptions struct {
		Colorize   bool
		ConfigPath string
		Env        []string
		EnvFiles   []string
		Verbose    bool
	}

	CommandRunner func(*config.Config) error
)

func Run(
	command string,
	config *config.Config,
	appOptions *AppOptions,
	cleanOptions *CleanOptions,
	runOptions *RunOptions,
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
