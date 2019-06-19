package paths

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ij-build/ij/logging"
)

type Copier struct {
	logger  logging.Logger
	project string
}

func NewCopier(logger logging.Logger, project string) *Copier {
	return &Copier{
		logger:  logger,
		project: project,
	}
}

func (c *Copier) Copy(src, dest string, blacklist map[string]struct{}) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	return c.copy(src, dest, blacklist, info, false)
}

func (c *Copier) copy(
	src string,
	dest string,
	blacklist map[string]struct{},
	info os.FileInfo,
	recursive bool,
) error {
	if _, ok := blacklist[src]; ok {
		c.logger.Debug(
			nil,
			"Skipping transfer of blacklisted file %s",
			c.displayPath(src),
		)

		return nil
	}

	if info.Mode()&os.ModeSymlink != 0 {
		c.logger.Debug(
			nil,
			"Skipping transfer of symlink %s",
			c.displayPath(src),
		)

		return nil
	}

	if !recursive {
		c.logger.Debug(
			nil,
			"Copying %s to %s",
			c.displayPath(src),
			c.displayPath(dest),
		)
	}

	if info.IsDir() {
		return c.copyAll(src, dest, blacklist, info)
	}

	return copyFile(src, dest, info)
}

func (c *Copier) copyAll(src, dest string, blacklist map[string]struct{}, info os.FileInfo) error {
	if err := EnsureDirExists(dest, info.Mode()); err != nil {
		return err
	}

	entries, err := DirContents(src)
	if err != nil {
		return err
	}

	for _, info := range entries {
		err := c.copy(
			filepath.Join(src, info.Name()),
			filepath.Join(dest, info.Name()),
			blacklist,
			info,
			true,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Copier) displayPath(path string) string {
	return fmt.Sprintf("~%s", path[len(c.project):])
}

//
// Helpers

func copyFile(src, dest string, info os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	if err := EnsureParentExists(dest, os.ModePerm); err != nil {
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if err := os.Chmod(destFile.Name(), info.Mode()); err != nil {
		return err
	}

	return nil
}
