package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func Dirnames(dirname string) ([]string, error) {
	dir, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	return dir.Readdirnames(0)
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

func BuildPath(parts ...string) (string, error) {
	fullPath := filepath.Join(parts...)

	if err := EnsureDirExists(filepath.Dir(fullPath)); err != nil {
		return "", err
	}

	return fullPath, nil
}
