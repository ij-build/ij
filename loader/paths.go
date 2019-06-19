package loader

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/ij-build/ij/paths"
)

var (
	defaultConfigPaths = []string{
		"ij.yaml",
		"ij.yml",
	}

	localOverridePaths = []string{
		"ij.override.yaml",
		"ij.override.yml",
	}

	globalOverridePaths = []string{
		filepath.Join(".ij", "override.yaml"),
	}
)

func GetConfigPath(path string) (string, error) {
	if path != "" {
		return path, nil
	}

	for _, path := range defaultConfigPaths {
		ok, err := paths.FileExists(path)
		if err != nil {
			return "", err
		}

		if ok {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not infer config file")
}

func getOverridePaths() ([]string, error) {
	found := []string{}

	current, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user (%s)", err.Error())
	}

	for _, path := range globalOverridePaths {
		path = filepath.Join(current.HomeDir, path)

		if ok, err := paths.FileExists(path); err == nil && ok {
			found = append(found, path)
			break
		}
	}

	for _, path := range localOverridePaths {
		if ok, err := paths.FileExists(path); err == nil && ok {
			found = append(found, path)
		}
	}

	return found, nil
}
