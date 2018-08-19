package paths

import "os"

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func EnsureDirExists(dirname string) error {
	exists, err := Exists(dirname)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return os.MkdirAll(dirname, os.ModeDir|os.ModePerm)
}
