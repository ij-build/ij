package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

func DirContents(dirname string) ([]os.FileInfo, error) {
	if exists, err := DirExists(dirname); err == nil && !exists {
		return nil, nil
	}

	dir, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	if info.IsDir() {
		return false, fmt.Errorf("%s exists but is not a file", path)
	}

	return true, nil
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s exists but is not a directory", path)
	}

	return true, nil
}

func EnsureDirExists(dirname string, mode os.FileMode) error {
	exists, err := DirExists(dirname)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return os.MkdirAll(dirname, mode|os.ModeDir)
}

func EnsureParentExists(path string, mode os.FileMode) error {
	return EnsureDirExists(filepath.Dir(path), mode)
}
