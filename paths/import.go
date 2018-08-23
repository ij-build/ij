package paths

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/logging"
)

type Importer struct {
	project      string
	scratch      string
	workspace    string
	importCopier *Copier
	exportCopier *Copier
}

func NewImporter(
	project string,
	scratch string,
	workspace string,
	blacklistPatterns []string,
	logger logging.Logger,
) (*Importer, error) {
	blacklist, err := constructBlacklist(project, blacklistPatterns)
	if err != nil {
		return nil, err
	}

	return &Importer{
		project:      project,
		scratch:      scratch,
		workspace:    workspace,
		importCopier: NewCopier(logger, project, blacklist),
		exportCopier: NewCopier(logger, project, map[string]struct{}{}),
	}, nil
}

// TODO - import to specific dest
// TODO - export to specific dest

func (i *Importer) Import(patterns []string) error {
	return runOnPatterns(patterns, "", i.importPath)
}

func (i *Importer) Export(patterns []string) error {
	return runOnPatterns(patterns, i.workspace, i.exportPath)
}

func (i *Importer) importPath(path string) error {
	return i.transferPath(
		path,
		"import",
		i.project,
		i.workspace,
		i.importCopier,
	)
}

func (i *Importer) exportPath(path string) error {
	return i.transferPath(
		path,
		"export",
		i.workspace,
		i.project,
		i.exportCopier,
	)
}

func (i *Importer) transferPath(
	path string,
	transferType string,
	srcRoot string,
	destRoot string,
	copier *Copier,
) error {
	src, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf(
			"failed to normalize %s path: %s",
			transferType,
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

	if strings.HasPrefix(path, srcRoot) {
		path = path[len(srcRoot):]
	}

	dest := filepath.Join(destRoot, path)

	if err := copier.Copy(src, dest); err != nil {
		return fmt.Errorf(
			"failed to %s path %s: %s",
			transferType,
			path,
			err.Error(),
		)
	}

	return nil
}
