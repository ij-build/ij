package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type (
	Config struct {
		Extends             string                     `json:"extends"`
		SSHIdentities       json.RawMessage            `json:"ssh-identities"`
		ForceSequential     bool                       `json:"force-sequential"`
		HealthcheckInterval Duration                   `json:"healthcheck-interval"`
		Registries          []json.RawMessage          `json:"registries"`
		Workspace           string                     `json:"workspace"`
		Environment         json.RawMessage            `json:"environment"`
		Import              *FileList                  `json:"import"`
		Export              *FileList                  `json:"export"`
		Tasks               map[string]json.RawMessage `json:"tasks"`
		Plans               map[string]*Plan           `json:"plans"`
		Metaplans           map[string][]string        `json:"metaplans"`
	}

	FileList struct {
		Files    json.RawMessage `json:"files"`
		Excludes json.RawMessage `json:"excludes"`
	}
)

func (c *Config) Translate(parent *config.Config) (*config.Config, error) {
	sshIdentities, err := unmarshalStringList(c.SSHIdentities)
	if err != nil {
		return nil, err
	}

	registries := []config.Registry{}
	for _, registry := range c.Registries {
		translated, err := translateRegistry(registry)
		if err != nil {
			return nil, err
		}

		registries = append(registries, translated)
	}

	environment, err := unmarshalStringList(c.Environment)
	if err != nil {
		return nil, err
	}

	importList, err := translateFileList(c.Import)
	if err != nil {
		return nil, err
	}

	exportList, err := translateFileList(c.Export)
	if err != nil {
		return nil, err
	}

	tasks := map[string]config.Task{}
	for name, task := range c.Tasks {
		translated, err := translateTask(parent, name, task)
		if err != nil {
			return nil, err
		}

		tasks[name] = translated
	}

	plans := map[string]*config.Plan{}
	for name, plan := range c.Plans {
		translated, err := plan.Translate(name)
		if err != nil {
			return nil, err
		}

		plans[name] = translated
	}

	return &config.Config{
		Extends:             c.Extends,
		SSHIdentities:       sshIdentities,
		ForceSequential:     c.ForceSequential,
		HealthcheckInterval: c.HealthcheckInterval.Duration,
		Registries:          registries,
		Workspace:           c.Workspace,
		Environment:         environment,
		Import:              importList,
		Export:              exportList,
		Tasks:               tasks,
		Plans:               plans,
		Metaplans:           c.Metaplans,
	}, nil
}

func translateFileList(list *FileList) (*config.FileList, error) {
	files, err := unmarshalStringList(list.Files)
	if err != nil {
		return nil, err
	}

	excludes, err := unmarshalStringList(list.Excludes)
	if err != nil {
		return nil, err
	}

	return &config.FileList{
		Files:    files,
		Excludes: excludes,
	}, nil
}
