package paths

import (
	"os"

	"github.com/efritz/ij/logging"
)

type Remover struct {
	project string
}

func NewRemover(project string) *Remover {
	return &Remover{
		project: project,
	}
}

func (r *Remover) Remove(patterns, blacklistPatterns []string, shouldRemove func(string) (bool, error)) error {
	blacklist, err := constructBlacklist(r.project, blacklistPatterns)
	if err != nil {
		return err
	}

	return runOnPatterns(patterns, r.project, logging.NilLogger, func(pair FilePair) error {
		dest, err := sanitize(pair.Dest, r.project)
		if err != nil {
			return err
		}

		if _, ok := blacklist[dest]; ok {
			return nil
		}

		// Ensure removal of a parent didn't also remove this file.
		if exists, err := PathExists(dest); err != nil || !exists {
			return err
		}

		if remove, err := shouldRemove(dest); err != nil || !remove {
			return err
		}

		return os.RemoveAll(dest)
	})
}
