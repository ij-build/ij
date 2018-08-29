package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Override struct {
	Registries    []json.RawMessage `json:"registries"`
	SSHIdentities json.RawMessage   `json:"ssh-identities"`
	Environment   []string          `json:"environment"`
	Excludes      json.RawMessage   `json:"exclude"`
}

func (c *Override) Translate() (*config.Override, error) {
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

	excludes, err := unmarshalStringList(c.Excludes)
	if err != nil {
		return nil, err
	}

	return &config.Override{
		SSHIdentities: sshIdentities,
		Environment:   c.Environment,
		Registries:    registries,
		Excludes:      excludes,
	}, nil
}
