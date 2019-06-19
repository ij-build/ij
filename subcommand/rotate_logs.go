package subcommand

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/options"
	"github.com/ij-build/ij/paths"
	"github.com/ij-build/ij/scratch"
)

func NewRotateLogsCommand(appOptions *options.AppOptions) CommandRunner {
	return func(config *config.Config) error {
		scratchPath := filepath.Join(
			appOptions.ScratchRoot,
			scratch.ScratchDir,
		)

		entries, err := paths.DirContents(scratchPath)
		if err != nil {
			return fmt.Errorf(
				"failed to read scratch directory: %s",
				err.Error(),
			)
		}

		if len(entries) == 0 {
			return nil
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].ModTime().After(entries[j].ModTime())
		})

		for _, info := range entries[1:] {
			if err := os.RemoveAll(filepath.Join(scratchPath, info.Name())); err != nil {
				return fmt.Errorf(
					"failed to delete run directory: %s",
					err.Error(),
				)
			}
		}

		return nil
	}
}
