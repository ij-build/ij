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
		project      string
		scratch      string
		workspace    string
		importCopier *Copier
		exportCopier *Copier
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
	blacklistPatterns []string,
	logger logging.Logger,
) (*Transferer, error) {
	blacklist, err := constructBlacklist(project, blacklistPatterns)
	if err != nil {
		return nil, err
	}

	return &Transferer{
		project:      project,
		scratch:      scratch,
		workspace:    workspace,
		importCopier: NewCopier(logger, project, blacklist),
		exportCopier: NewCopier(logger, project, map[string]struct{}{}),
	}, nil
}

func (t *Transferer) Import(patterns []string) error {
	return runOnPatterns(patterns, t.project, true, true, t.importPath)
}

func (t *Transferer) Export(patterns []string) error {
	return runOnPatterns(patterns, t.workspace, true, true, t.exportPath)
}

func (t *Transferer) importPath(pair filePair) error {
	return t.transferPath(
		pair.src,
		pair.dest,
		"import",
		t.project,
		t.workspace,
		t.importCopier,
	)
}

func (t *Transferer) exportPath(pair filePair) error {
	return t.transferPath(
		pair.src,
		pair.dest,
		"export",
		t.workspace,
		t.project,
		t.exportCopier,
	)
}

func (t *Transferer) transferPath(
	rawSrc string,
	rawDest string,
	transferType string,
	srcRoot string,
	destRoot string,
	copier *Copier,
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

	if err := copier.Copy(src, dest); err != nil {
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

func constructBlacklist(project string, patterns []string) (map[string]struct{}, error) {
	var (
		blacklist   = map[string]struct{}{}
		allPatterns = append(DefaultBlacklist, patterns...)
	)

	err := runOnPatterns(allPatterns, project, false, false, func(pair filePair) error {
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
	strict bool,
	target func(filePair) error,
) error {
	for _, pattern := range patterns {
		if err := runOnPattern(pattern, root, split, strict, target); err != nil {
			return err
		}
	}

	return nil
}

func runOnPattern(
	pattern string,
	root string,
	split bool,
	strict bool,
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

	return runOnGlobPattern(pattern, root, strict, target)
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
	strict bool,
	target func(filePair) error,
) error {
	paths, err := zglob.Glob(filepath.Join(root, pattern))
	if err != nil {
		if err == os.ErrNotExist && !strict {
			return nil
		}

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
