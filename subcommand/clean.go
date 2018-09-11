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

func NewCleanCommand(appOptions *options.AppOptions, cleanOptions *options.CleanOptions) CommandRunner {
	return func(config *config.Config) error {
		err := paths.NewRemover(appOptions.ProjectDir).Remove(
			config.Export.Files,
			config.Export.CleanExcludes,
			cleanPromptFactory(appOptions.ProjectDir, cleanOptions.ForceClean),
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

func cleanPromptFactory(workingDir string, force bool) func(string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	return func(path string) (bool, error) {
		if force {
			return true, nil
		}

		fmt.Printf("remove '%s'? ", path[len(workingDir):])

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
