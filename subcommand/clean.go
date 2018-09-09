package subcommand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/options"
	"github.com/efritz/ij/paths"
)

func NewCleanCommand(cleanOptions *options.CleanOptions) CommandRunner {
	return func(config *config.Config) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"failed to get working directory: %s",
				err.Error(),
			)
		}

		err = paths.NewRemover(wd).Remove(
			config.Export.Files,
			config.Export.CleanExcludes,
			cleanPromptFactory(wd, cleanOptions.ForceClean),
		)

		if err != nil {
			return fmt.Errorf(
				"Failed to clean exported files: %s",
				err.Error(),
			)
		}

		return nil
	}
}

func cleanPromptFactory(wd string, force bool) func(string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	return func(path string) (bool, error) {
		if force {
			return true, nil
		}

		fmt.Printf("remove '%s'? ", path[len(wd):])

		text, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		if strings.ToLower(strings.TrimSpace(text)) == "y" {
			return true, nil
		}

		return false, nil
	}
}
