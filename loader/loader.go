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
		loadedConfigs     map[string]*jsonconfig.Config
		loadedOverrides   map[string]*config.Override
		dependencyGraph   *topsort.Graph
		pathSubstitutions map[string]string
	}

	jsonEnvelope struct {
		Tasks     map[string]json.RawMessage `json:"tasks"`
		Plans     map[string]json.RawMessage `json:"plans"`
		Metaplans map[string]json.RawMessage `json:"metaplans"`
	}
)

func NewLoader() *Loader {
	return &Loader{
		loadedConfigs:     map[string]*jsonconfig.Config{},
		loadedOverrides:   map[string]*config.Override{},
		dependencyGraph:   topsort.NewGraph(),
		pathSubstitutions: map[string]string{},
	}
}

func (l *Loader) LoadPathSubstitutions(overridePaths []string) error {
	for _, path := range overridePaths {
		override, err := l.readOverride(path)
		if err != nil {
			return err
		}

		for k, v := range override.Options.PathSubstitutions {
			l.pathSubstitutions[k] = v
		}
	}

	return nil
}

func (l *Loader) Load(path string) (*config.Config, error) {
	path = l.normalizePath(path, "")

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
		child, err := l.loadedConfigs[path].Translate(config)
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

func (l *Loader) ApplyOverrides(config *config.Config, overridePaths []string) error {
	for _, path := range overridePaths {
		if err := l.applyOverride(config, path); err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) readConfigs(path string) error {
	if _, ok := l.loadedConfigs[path]; ok {
		return nil
	}

	config, err := l.readConfig(path)
	if err != nil {
		return err
	}

	l.loadedConfigs[path] = config
	l.dependencyGraph.AddNode(path)

	extends, err := util.UnmarshalStringList(config.Extends)
	if err != nil {
		return err
	}

	for _, parent := range extends {
		parentPath := l.normalizePath(parent, path)

		if err := l.readConfigs(parentPath); err != nil {
			return err
		}

		l.dependencyGraph.AddEdge(path, parentPath)
	}

	for i := 1; i < len(extends); i++ {
		path1 := l.normalizePath(extends[i], path)
		path2 := l.normalizePath(extends[i-1], path)

		l.dependencyGraph.AddEdge(path1, path2)
	}

	return nil
}

func (l *Loader) readConfig(path string) (*jsonconfig.Config, error) {
	data, err := readPath(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load config %s: %s",
			path,
			err.Error(),
		)
	}

	if err := validateConfig(path, data); err != nil {
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
	override, err := l.readOverride(path)
	if err != nil {
		return err
	}

	config.ApplyOverride(override)
	return nil
}

func (l *Loader) readOverride(path string) (*config.Override, error) {
	if override, ok := l.loadedOverrides[path]; ok {
		return override, nil
	}

	data, err := readPath(path)
	if err != nil {
		return nil, err
	}

	if err := validateOverride(path, data); err != nil {
		return nil, err
	}

	payload := &jsonconfig.Override{
		Options: &jsonconfig.Options{},
		Import:  &jsonconfig.ImportFileList{},
		Export:  &jsonconfig.ExportFileList{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	override, err := payload.Translate()
	if err != nil {
		return nil, err
	}

	l.loadedOverrides[path] = override
	return override, nil
}

func (l *Loader) normalizePath(path, source string) string {
	rawPath := buildPath(path, source)

	realPath := rawPath
	for k, v := range l.pathSubstitutions {
		realPath = strings.Replace(realPath, k, v, -1)
	}

	// Transformed to a differing relative path
	if rawPath != realPath && !isURL(realPath) && !filepath.IsAbs(realPath) {
		realPath = buildPath(realPath, source)
	}

	return realPath
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

	if err := loader.LoadPathSubstitutions(overridePaths); err != nil {
		return nil, fmt.Errorf(
			"failed to load path substitutions from overrride file: %s",
			err.Error(),
		)
	}

	cfg, err := loader.Load(path)
	if err != nil {
		return nil, err
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

		fileLines, err := environment.NormalizeEnvironmentFile(string(content))
		if err != nil {
			return nil, err
		}

		lines = append(lines, fileLines...)
	}

	return lines, nil
}

func buildPath(path, source string) string {
	if isURL(path) || isURL(source) {
		return path
	}

	if source == "" {
		return path
	}

	return filepath.Join(filepath.Dir(source), path)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from remote server", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func validateConfig(path string, data []byte) error {
	if err := schema.Validate("schema/config.yaml", data); err != nil {
		return fmt.Errorf("failed to validate config %s: %s", path, err.Error())
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

func validateOverride(path string, data []byte) error {
	if err := schema.Validate("schema/override.yaml", data); err != nil {
		return fmt.Errorf("failed to validate override %s: %s", path, err.Error())
	}

	return nil
}
