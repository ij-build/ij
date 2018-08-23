package paths

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

type filePair struct {
	src  string
	dest string
}

func runOnPatterns(
	patterns []string,
	root string,
	split bool,
	target func(filePair) error,
) error {
	for _, pattern := range patterns {
		if err := runOnPattern(pattern, root, split, target); err != nil {
			return err
		}
	}

	return nil
}

func runOnPattern(
	pattern string,
	root string,
	split bool,
	target func(filePair) error,
) error {
	if strings.Contains(pattern, ":") {
		if strings.Contains(pattern, "*") {
			return fmt.Errorf(
				"illegal pattern %s (wildcards and explicit destinations are mutually exclusive)",
				pattern,
			)
		}

		return runOnSplitPattern(pattern, root, target)
	}

	return runOnGlobPattern(pattern, root, target)
}

func runOnSplitPattern(
	pattern string,
	root string,
	target func(filePair) error,
) error {
	src, dest := splitPath(pattern)

	return target(filePair{
		src:  filepath.Join(root, src),
		dest: filepath.Join(root, dest),
	})
}

func runOnGlobPattern(
	pattern string,
	root string,
	target func(filePair) error,
) error {
	paths, err := zglob.Glob(filepath.Join(root, pattern))
	if err != nil {
		return fmt.Errorf(
			"failed to glob pattern %s: %s",
			pattern,
			err.Error(),
		)
	}

	for _, path := range paths {
		if err := target(filePair{path, path}); err != nil {
			return err
		}
	}

	return nil
}

func splitPath(path string) (string, string) {
	if parts := strings.SplitN(path, ":", 2); len(parts) == 2 {
		return parts[0], parts[1]
	}

	return path, path
}
