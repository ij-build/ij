package paths

import (
	"fmt"
	"path/filepath"

	"github.com/ij-build/ij/logging"
)

type Transferer struct {
	project   string
	scratch   string
	workspace string
	logger    logging.Logger
	copier    *Copier
}

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

	return runOnPatterns(patterns, t.project, t.logger, func(pair FilePair) error {
		return t.transferPath(pair.Src, pair.Dest, t.project, t.workspace, blacklist)
	})
}

func (t *Transferer) Export(patterns, blacklistPatterns []string) error {
	blacklist, err := constructBlacklist(t.workspace, blacklistPatterns)
	if err != nil {
		return err
	}

	return runOnPatterns(patterns, t.workspace, t.logger, func(pair FilePair) error {
		return t.transferPath(pair.Src, pair.Dest, t.workspace, t.project, blacklist)
	})
}

func (t *Transferer) transferPath(
	rawSrc string,
	rawDest string,
	srcRoot string,
	destRoot string,
	blacklist map[string]struct{},
) error {
	src, err := sanitize(rawSrc, srcRoot)
	if err != nil {
		return err
	}

	dest := filepath.Join(destRoot, rawDest[len(srcRoot):])

	if err := t.copier.Copy(src, dest, blacklist); err != nil {
		return fmt.Errorf(
			"failed to copy path '%s': %s",
			rawSrc,
			err.Error(),
		)
	}

	return nil
}
