package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Config struct {
	Extends       string                     `json:"extends"`
	SSHIdentities json.RawMessage            `json:"ssh-identities"`
	Workspace     string                     `json:"workspace"`
	Environment   []string                   `json:"environment"`
	Tasks         map[string]json.RawMessage `json:"tasks"`
	Plans         map[string]*Plan           `json:"plans"`
	Metaplans     map[string][]string        `json:"metaplans"`
	Imports       json.RawMessage            `json:"import"`
	Exports       json.RawMessage            `json:"export"`
	Excludes      json.RawMessage            `json:"exclude"`
}

func (c *Config) Translate() (*config.Config, error) {
	sshIdentities, err := unmarshalStringList(c.SSHIdentities)
	if err != nil {
		return nil, err
	}

	tasks := map[string]config.Task{}
	for name, task := range c.Tasks {
		translated, err := translateTask(name, task)
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

	return &config.Config{
		Extends:       c.Extends,
		SSHIdentities: sshIdentities,
		Workspace:     c.Workspace,
		Environment:   c.Environment,
		Tasks:         tasks,
		Plans:         plans,
		Metaplans:     c.Metaplans,
		Imports:       imports,
		Exports:       exports,
		Excludes:      excludes,
	}, nil
}
