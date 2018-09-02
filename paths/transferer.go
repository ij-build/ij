package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/logging"
	zglob "github.com/mattn/go-zglob"
)

type (
	Transferer struct {
		project   string
		scratch   string
		workspace string
		logger    logging.Logger
		copier    *Copier
	}

	filePair struct {
		src  string
		dest string
	}
)

var DefaultBlacklist = []string{
	".ij",
	".git",
}

func NewTransferer(
	project string,
	scratch string,
	workspace string,
	logger logging.Logger,
) *Transferer {
	return &Transferer{
		project:   project,
		scratch:   scratch,
		workspace: workspace,
		logger:    logger,
		copier:    NewCopier(logger, project),
	}
}

func (t *Transferer) Import(patterns, blacklistPatterns []string) error {
	blacklist, err := constructBlacklist(t.project, blacklistPatterns)
	if err != nil {
		return err
	}

	return runOnPatterns(patterns, t.project, true, t.logger, func(pair filePair) error {
		return t.transferPath(pair.src, pair.dest, t.project, t.workspace, blacklist, "import")
	})
}

func (t *Transferer) Export(patterns, blacklistPatterns []string) error {
	blacklist, err := constructBlacklist(t.workspace, blacklistPatterns)
	if err != nil {
		return err
	}

	return runOnPatterns(patterns, t.workspace, true, t.logger, func(pair filePair) error {
		return t.transferPath(pair.src, pair.dest, t.workspace, t.project, blacklist, "export")
	})
}

func (t *Transferer) transferPath(
	rawSrc string,
	rawDest string,
	srcRoot string,
	destRoot string,
	blacklist map[string]struct{},
	transferType string,
) error {
	src, err := filepath.Abs(rawSrc)
	if err != nil {
		return fmt.Errorf(
			"failed to normalize %s path %s: %s",
			transferType,
			rawSrc,
			err.Error(),
		)
	}

	if !strings.HasPrefix(src, srcRoot) {
		return fmt.Errorf(
			"%s file is outside of workspace directory: %s",
			transferType,
			src,
		)
	}

	dest := filepath.Join(destRoot, rawDest[len(srcRoot):])

	if err := t.copier.Copy(src, dest, blacklist); err != nil {
		return fmt.Errorf(
			"failed to %s path %s: %s",
			transferType,
			rawSrc,
			err.Error(),
		)
	}

	return nil
}

//
// Helpers

func constructBlacklist(root string, patterns []string) (map[string]struct{}, error) {
	var (
		blacklist   = map[string]struct{}{}
		allPatterns = append(DefaultBlacklist, patterns...)
	)

	err := runOnPatterns(allPatterns, root, false, logging.NilLogger, func(pair filePair) error {
		blacklist[pair.src] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blacklist, nil
}

func runOnPatterns(
	patterns []string,
	root string,
	split bool,
	logger logging.Logger,
	target func(filePair) error,
) error {
	for _, pattern := range patterns {
		if err := runOnPattern(pattern, root, split, logger, target); err != nil {
			return err
		}
	}

	return nil
}

func runOnPattern(
	pattern string,
	root string,
	split bool,
	logger logging.Logger,
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

	return runOnGlobPattern(pattern, root, logger, target)
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
	logger logging.Logger,
	target func(filePair) error,
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
