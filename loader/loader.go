package loader

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
)

type Loader struct {
	loaded map[string]struct{}
}

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

	if err := validateWithSchema(data); err != nil {
		return nil, err
	}

	child := &config.Config{}
	if err := json.Unmarshal(data, child); err != nil {
		return nil, err
	}

	if child.Tasks == nil {
		child.Tasks = map[string]*config.Task{}
	}

	if child.Plans == nil {
		child.Plans = map[string]*config.Plan{}
	}

	if err := unmarshalFileList(child); err != nil {
		return nil, err
	}

	if err := unmarshalStageTasks(child); err != nil {
		return nil, err
	}

	populateTaskNames(child)
	populatePlanNames(child)

	return l.resolveParent(child)
}

func (l *Loader) resolveParent(config *config.Config) (*config.Config, error) {
	if config.Extends == "" {
		return config, nil
	}

	if _, ok := l.loaded[config.Extends]; ok {
		return nil, fmt.Errorf(
			"failed to extend from %s (extension is cyclic)",
			config.Extends,
		)
	}

	l.loaded[config.Extends] = struct{}{}

	parent, err := l.Load(config.Extends)
	if err != nil {
		return nil, err
	}

	if err := mergeConfigs(config, parent); err != nil {
		return nil, err
	}

	return parent, nil
}
