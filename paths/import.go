package paths

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/logging"
)

type Transferer struct {
	project      string
	scratch      string
	workspace    string
	importCopier *Copier
	exportCopier *Copier
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
	return runOnPatterns(patterns, "", true, t.importPath)
}

func (t *Transferer) Export(patterns []string) error {
	return runOnPatterns(patterns, t.workspace, true, t.exportPath)
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

	if strings.HasPrefix(rawDest, srcRoot) {
		rawDest = rawDest[len(srcRoot):]
	}

	dest := filepath.Join(destRoot, rawDest)

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
