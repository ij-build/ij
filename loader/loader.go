package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/environment"
	"github.com/efritz/ij/util"
	"github.com/ghodss/yaml"
	"github.com/stevenle/topsort"

	"github.com/efritz/ij/loader/jsonconfig"
	"github.com/efritz/ij/loader/schema"
)

type (
	Loader struct {
		loaded          map[string]*jsonconfig.Config
		dependencyGraph *topsort.Graph
	}

	jsonEnvelope struct {
		Tasks     map[string]json.RawMessage `json:"tasks"`
		Plans     map[string]json.RawMessage `json:"plans"`
		Metaplans map[string]json.RawMessage `json:"metaplans"`
	}
)

func NewLoader() *Loader {
	return &Loader{
		loaded:          map[string]*jsonconfig.Config{},
		dependencyGraph: topsort.NewGraph(),
	}
}

func (l *Loader) Load(path string) (*config.Config, error) {
	if err := l.readConfigs(path); err != nil {
		return nil, err
	}

	order, err := l.dependencyGraph.TopSort(path)
	if err != nil {
		// Error messages starts with "Cycle error: "
		return nil, fmt.Errorf("failed to extend cyclic config (%s)", err.Error()[13:])
	}

	var config *config.Config
	for _, path := range order {
		child, err := l.loaded[path].Translate(config)
		if err != nil {
			return nil, err
		}

		if config == nil {
			config = child
		} else {
			if err := config.Merge(child); err != nil {
				return nil, err
			}
		}
	}

	return config, nil
}

func (l *Loader) ApplyOverrides(config *config.Config, paths []string) error {
	for _, path := range paths {
		if err := l.applyOverride(config, path); err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) readConfigs(path string) error {
	if _, ok := l.loaded[path]; ok {
		return nil
	}

	config, err := l.readConfig(path)
	if err != nil {
		return err
	}

	l.loaded[path] = config
	l.dependencyGraph.AddNode(path)

	extends, err := util.UnmarshalStringList(config.Extends)
	if err != nil {
		return err
	}

	for _, parent := range extends {
		parentPath := buildPath(parent, path)

		if err := l.readConfigs(parentPath); err != nil {
			return err
		}

		l.dependencyGraph.AddEdge(path, parentPath)
	}

	for i := 1; i < len(extends); i++ {
		path1 := buildPath(extends[i], path)
		path2 := buildPath(extends[i-1], path)

		l.dependencyGraph.AddEdge(path1, path2)
	}

	return nil
}

func (l *Loader) readConfig(path string) (*jsonconfig.Config, error) {
	data, err := readPath(path)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(data); err != nil {
		return nil, err
	}

	payload := &jsonconfig.Config{
		Options:   &jsonconfig.Options{},
		Import:    &jsonconfig.ImportFileList{},
		Export:    &jsonconfig.ExportFileList{},
		Tasks:     map[string]json.RawMessage{},
		Plans:     map[string]*jsonconfig.Plan{},
		Metaplans: map[string][]string{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (l *Loader) applyOverride(config *config.Config, path string) error {
	data, err := readPath(path)
	if err != nil {
		return err
	}

	if err := schema.Validate("schema/override.yaml", data); err != nil {
		return fmt.Errorf("failed to validate override file: %s", err.Error())
	}

	payload := &jsonconfig.Override{
		Options: &jsonconfig.Options{},
		Import:  &jsonconfig.ImportFileList{},
		Export:  &jsonconfig.ExportFileList{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return err
	}

	override, err := payload.Translate()
	if err != nil {
		return err
	}

	config.ApplyOverride(override)
	return nil
}

//
// Command Line

func LoadFile(path string, override *config.Override) (*config.Config, error) {
	overridePaths, err := getOverridePaths()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to determine override paths: %s",
			err.Error(),
		)
	}

	loader := NewLoader()

	cfg, err := loader.Load(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load config: %s",
			err.Error(),
		)
	}

	if err := loader.ApplyOverrides(cfg, overridePaths); err != nil {
		return nil, fmt.Errorf(
			"failed to apply overrides: %s",
			err.Error(),
		)
	}

	if override != nil {
		cfg.ApplyOverride(override)
	}

	envFromFile, err := applyEnvironmentFiles(cfg.EnvironmentFiles)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read environment file: %s",
			err.Error(),
		)
	}

	cfg.Environment = append(
		environment.Default().Serialize(),
		append(
			cfg.Environment,
			envFromFile...,
		)...,
	)

	if err := cfg.Resolve(); err != nil {
		return nil, fmt.Errorf(
			"failed to resolve config: %s",
			err.Error(),
		)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf(
			"failed to validate config: %s",
			err.Error(),
		)
	}

	return cfg, nil
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

//
// Helpers

func readPath(path string) ([]byte, error) {
	data, err := chooseReader(path)(path)
	if err != nil {
		return nil, err
	}

	return yaml.YAMLToJSON(data)
}

func chooseReader(path string) func(string) ([]byte, error) {
	if isURL(path) {
		return readRemoteFile
	}

	return ioutil.ReadFile
}

func readRemoteFile(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func buildPath(path, source string) string {
	if isURL(path) || isURL(source) {
		return path
	}

	return filepath.Join(filepath.Dir(source), path)
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func validateConfig(data []byte) error {
	if err := schema.Validate("schema/config.yaml", data); err != nil {
		return fmt.Errorf("failed to validate config: %s", err.Error())
	}

	payload := &jsonEnvelope{
		Tasks:     map[string]json.RawMessage{},
		Plans:     map[string]json.RawMessage{},
		Metaplans: map[string]json.RawMessage{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return err
	}

	for name, plan := range payload.Plans {
		if err := schema.Validate("schema/plan.yaml", plan); err != nil {
			return fmt.Errorf("failed to validate plan %s: %s", name, err.Error())
		}
	}

	for name, metaplan := range payload.Metaplans {
		if err := schema.Validate("schema/metaplan.yaml", metaplan); err != nil {
			return fmt.Errorf("failed to validate metaplan %s: %s", name, err.Error())
		}
	}

	return nil
}
