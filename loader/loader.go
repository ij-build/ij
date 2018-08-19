package loader

import (
	"encoding/json"

	"github.com/efritz/pvc/config"
)

func LoadFile(path string) (*config.Config, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return Load(data)
}

// TODO - laod url and other stuff

func Load(data []byte) (*config.Config, error) {
	if err := validateWithSchema(data); err != nil {
		return nil, err
	}

	cfg := &config.Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	for name, task := range cfg.Tasks {
		task.Name = name
	}

	for name, plan := range cfg.Plans {
		plan.Name = name
	}

	// TODO - do other validation
	return cfg, nil
}
