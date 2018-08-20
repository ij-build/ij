package loader

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/pvc/config"
)

type Loader struct {
	loaded []string
}

func NewLoader() *Loader {
	return &Loader{
		loaded: []string{},
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

	for _, path := range l.loaded {
		if path == config.Extends {
			return nil, fmt.Errorf(
				"failed to extend from %s (extension is cyclic)",
				config.Extends,
			)
		}
	}

	l.loaded = append(l.loaded, config.Extends)

	parent, err := l.Load(config.Extends)
	if err != nil {
		return nil, err
	}

	if err := merge(config, parent); err != nil {
		return nil, err
	}

	return parent, nil
}
