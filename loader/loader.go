package loader

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/pvc/config"
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

	config := &config.Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if err := unmarshalStageTasks(config); err != nil {
		return nil, err
	}

	populateTaskNames(config)
	populatePlanNames(config)

	return l.resolveParent(config)
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

	if err := merge(config, parent); err != nil {
		return nil, err
	}

	return parent, nil
}
