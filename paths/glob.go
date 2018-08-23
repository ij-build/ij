package paths

import (
	"fmt"
	"path/filepath"

	"github.com/mattn/go-zglob"
)

func runOnPatterns(patterns []string, root string, f func(string) error) error {
	for _, pattern := range patterns {
		paths, err := zglob.Glob(filepath.Join(root, pattern))
		if err != nil {
			return fmt.Errorf(
				"failed to glob pattern %s: %s",
				pattern,
				err.Error(),
			)
		}

		for _, path := range paths {
			if err := f(path); err != nil {
				return err
			}
		}
	}

	return nil
}
