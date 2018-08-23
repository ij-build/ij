package paths

import (
	"io"
	"os"
)

func copyFile(src, dest string, info os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	if err := EnsureParentExists(dest, os.ModePerm); err != nil {
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
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
