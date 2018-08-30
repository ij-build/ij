package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Override struct {
	SSHIdentities       json.RawMessage   `json:"ssh-identities"`
	ForceSequential     bool              `json:"force-sequential"`
	HealthcheckInterval Duration          `json:"healthcheck-interval"`
	Registries          []json.RawMessage `json:"registries"`
	Environment         json.RawMessage   `json:"environment"`
	Import              *FileList         `json:"import"`
	Export              *FileList         `json:"export"`
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

	return &config.Override{
		SSHIdentities:       sshIdentities,
		ForceSequential:     c.ForceSequential,
		HealthcheckInterval: c.HealthcheckInterval.Duration,
		Registries:          registries,
		Environment:         environment,
		ImportExcludes:      importList.Excludes,
		ExportExcludes:      exportList.Excludes,
	}, nil
}
