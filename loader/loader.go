package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/efritz/ij/config"
	"github.com/ghodss/yaml"

	"github.com/efritz/ij/loader/jsonconfig"
	"github.com/efritz/ij/loader/schema"
)

type (
	Loader struct {
		loaded map[string]struct{}
	}

	jsonEnvelope struct {
		Tasks     map[string]json.RawMessage `json:"tasks"`
		Plans     map[string]json.RawMessage `json:"plans"`
		Metaplans map[string]json.RawMessage `json:"metaplans"`
	}
)

func NewLoader() *Loader {
	return &Loader{
		loaded: map[string]struct{}{},
	}
}

func (l *Loader) Load(path string) (*config.Config, error) {
	data, err := readPath(path)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(data); err != nil {
		return nil, err
	}

	payload := &jsonconfig.Config{
		Options:   &jsonconfig.Options{},
		Import:    &jsonconfig.FileList{},
		Export:    &jsonconfig.FileList{},
		Tasks:     map[string]json.RawMessage{},
		Plans:     map[string]*jsonconfig.Plan{},
		Metaplans: map[string][]string{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	return l.resolveParent(payload, path)
}

func (l *Loader) ApplyOverrides(config *config.Config, paths []string) error {
	for _, path := range paths {
		if err := l.applyOverride(config, path); err != nil {
			return err
		}
	}

	return nil
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
		Import:  &jsonconfig.FileList{},
		Export:  &jsonconfig.FileList{},
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

func (l *Loader) resolveParent(
	config *jsonconfig.Config,
	childPath string,
) (*config.Config, error) {
	if config.Extends == "" {
		child, err := config.Translate(nil)
		if err != nil {
			return nil, err
		}

		return child, nil
	}

	if _, ok := l.loaded[config.Extends]; ok {
		return nil, fmt.Errorf(
			"failed to extend from %s (extension is cyclic)",
			config.Extends,
		)
	}

	l.loaded[config.Extends] = struct{}{}

	parent, err := l.Load(buildPath(config.Extends, childPath))
	if err != nil {
		return nil, err
	}

	child, err := config.Translate(parent)
	if err != nil {
		return nil, err
	}

	if err := parent.Merge(child); err != nil {
		return nil, err
	}

	return parent, nil
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
