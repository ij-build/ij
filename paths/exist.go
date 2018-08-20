package paths

import (
	"fmt"
	"os"
)

func EnsureDirExists(dirname string) error {
	exists, err := DirExists(dirname)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return os.MkdirAll(dirname, os.ModeDir|os.ModePerm)
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
