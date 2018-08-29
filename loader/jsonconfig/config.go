package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Config struct {
	Extends       string                     `json:"extends"`
	Registries    []json.RawMessage          `json:"registries"`
	SSHIdentities json.RawMessage            `json:"ssh-identities"`
	Environment   []string                   `json:"environment"`
	Imports       json.RawMessage            `json:"import"`
	Exports       json.RawMessage            `json:"export"`
	Excludes      json.RawMessage            `json:"exclude"`
	Workspace     string                     `json:"workspace"`
	Tasks         map[string]json.RawMessage `json:"tasks"`
	Plans         map[string]*Plan           `json:"plans"`
	Metaplans     map[string][]string        `json:"metaplans"`
}

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

	imports, err := unmarshalStringList(c.Imports)
	if err != nil {
		return nil, err
	}

	exports, err := unmarshalStringList(c.Exports)
	if err != nil {
		return nil, err
	}

	excludes, err := unmarshalStringList(c.Excludes)
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
		Extends:       c.Extends,
		SSHIdentities: sshIdentities,
		Workspace:     c.Workspace,
		Environment:   c.Environment,
		Registries:    registries,
		Imports:       imports,
		Exports:       exports,
		Excludes:      excludes,
		Tasks:         tasks,
		Plans:         plans,
		Metaplans:     c.Metaplans,
	}, nil
}
