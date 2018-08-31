package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/loader"
	"github.com/efritz/ij/logging"
	"github.com/efritz/ij/paths"
)

var (
	localOverridePaths = []string{
		"ij.override.yaml",
		"ij.override.yml",
	}

	globalOverridePaths = []string{
		filepath.Join(".ij", "override.yaml"),
	}
)

func loadConfig() (*config.Config, bool) {
	overridePaths, err := getOverridePaths()
	if err != nil {
		logging.EmergencyLog(
			"error: failed to determine override paths: %s",
			err.Error(),
		)

		return nil, false
	}

	loader := loader.NewLoader()

	cfg, err := loader.Load(*configPath)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to load config: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := loader.ApplyOverrides(cfg, overridePaths); err != nil {
		logging.EmergencyLog(
			"error: failed to apply overrides: %s",
			err.Error(),
		)

		return nil, false
	}

	cfg.ApplyOverride(&config.Override{
		Options: &config.Options{
			SSHIdentities:       *sshIdentities,
			ForceSequential:     *forceSequential,
			HealthcheckInterval: *healthcheckInterval,
		},
		EnvironmentFiles: *envFiles,
	})

	envFromFile, err := applyEnvironmentFiles(cfg.EnvironmentFiles)
	if err != nil {
		logging.EmergencyLog(
			"error: failed to read environment file: %s",
			err.Error(),
		)

		return nil, false
	}

	cfg.Environment = append(
		environment.Default().Serialize(),
		append(
			cfg.Environment,
			envFromFile...,
		)...,
	)

	if err := cfg.Resolve(); err != nil {
		logging.EmergencyLog(
			"error: failed to resolve config: %s",
			err.Error(),
		)

		return nil, false
	}

	if err := cfg.Validate(); err != nil {
		logging.EmergencyLog(
			"error: failed to validate config: %s",
			err.Error(),
		)

		return nil, false
	}

	for _, name := range *plans {
		if !cfg.IsPlanDefined(name) {
			logging.EmergencyLog(
				"error: unknown plan %s",
				name,
			)

			return nil, false
		}
	}

	return cfg, true
}

func getOverridePaths() ([]string, error) {
	current, err := user.Current()
	if err != nil {
		return nil, err
	}

	found := []string{}
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

func applyEnvironmentFiles(environmentFiles []string) ([]string, error) {
	lines := []string{}
	for _, path := range environmentFiles {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		lines, err := environment.NormalizeEnvironmentFile(string(content))
		if err != nil {
			return nil, err
		}

		lines = append(lines, lines...)
	}

	return lines, nil
}
