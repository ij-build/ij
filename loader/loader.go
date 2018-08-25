package loader

import (
	"encoding/json"
	"fmt"

	"github.com/efritz/ij/config"
	"github.com/efritz/ij/loader/jsonconfig"
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

	payload := &jsonconfig.Config{
		Tasks:     map[string]json.RawMessage{},
		Plans:     map[string]*jsonconfig.Plan{},
		Metaplans: map[string][]string{},
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	child, err := payload.Translate()
	if err != nil {
		return nil, err
	}

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

	if err := parent.Merge(config); err != nil {
		return nil, err
	}

	return parent, nil
}
