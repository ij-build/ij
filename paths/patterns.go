package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ij-build/ij/logging"
	"github.com/mattn/go-zglob"
)

type FilePair struct {
	Src  string
	Dest string
}

func runOnPatterns(
	patterns []string,
	root string,
	logger logging.Logger,
	target func(FilePair) error,
) error {
	for _, pattern := range patterns {
		if err := runOnPattern(pattern, root, logger, target); err != nil {
			return err
		}
	}

	return nil
}

func runOnPattern(
	pattern string,
	root string,
	logger logging.Logger,
	target func(FilePair) error,
) error {
	if strings.Contains(pattern, ":") {
		if strings.Contains(pattern, "*") {
			return fmt.Errorf(
				"illegal pattern %s (wildcards and explicit destinations are mutually exclusive)",
				pattern,
			)
		}

		return runOnSplitPattern(pattern, root, logger, target)
	}

	return runOnGlobPattern(pattern, root, logger, target)
}

func sanitize(path, root string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf(
			"failed to normalize path '%s': %s",
			path,
			err.Error(),
		)
	}

	if !strings.HasPrefix(absPath, root) {
		return "", fmt.Errorf(
			"path '%s' is outside of workspace directory",
			absPath,
		)
	}

	return absPath, nil
}

//
// Helpers

func constructBlacklist(root string, patterns []string) (map[string]struct{}, error) {
	var (
		blacklist   = map[string]struct{}{}
		allPatterns = append(DefaultBlacklist, patterns...)
	)

	for _, pattern := range allPatterns {
		if strings.Contains(pattern, ":") {
			return nil, fmt.Errorf("blacklist contains destination path: %s", pattern)
		}
	}

	err := runOnPatterns(allPatterns, root, logging.NilLogger, func(pair FilePair) error {
		blacklist[pair.Src] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blacklist, nil
}

func runOnSplitPattern(
	pattern string,
	root string,
	logger logging.Logger,
	target func(FilePair) error,
) error {
	src, dest := splitPath(pattern)

	if exists, err := PathExists(filepath.Join(root, src)); err != nil || !exists {
		logger.Warn(
			nil,
			"no files matched the pattern '%s'",
			src,
		)

		return err
	}

	return target(FilePair{
		Src:  filepath.Join(root, src),
		Dest: filepath.Join(root, dest),
	})
}

func runOnGlobPattern(
	pattern string,
	root string,
	logger logging.Logger,
	target func(FilePair) error,
) error {
	paths, err := zglob.Glob(filepath.Join(root, pattern))
	if err != nil {
		if err == os.ErrNotExist {
			logger.Warn(
				nil,
				"no files matched the pattern '%s'",
				pattern,
			)

			return nil
		}

		return err
	}

	if len(paths) == 0 {
		logger.Warn(
			nil,
			"no files matched the pattern '%s'",
			pattern,
		)
	}

	for _, path := range paths {
		if err := target(FilePair{path, path}); err != nil {
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
