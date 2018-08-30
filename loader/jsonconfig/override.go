package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Override struct {
	Options     *Options          `json:"options"`
	Registries  []json.RawMessage `json:"registries"`
	Environment json.RawMessage   `json:"environment"`
	Import      *FileList         `json:"import"`
	Export      *FileList         `json:"export"`
}

func (o *Override) Translate() (*config.Override, error) {
	options, err := o.Options.Translate()
	if err != nil {
		return nil, err
	}

	registries := []config.Registry{}
	for _, registry := range o.Registries {
		translated, err := translateRegistry(registry)
		if err != nil {
			return nil, err
		}

		registries = append(registries, translated)
	}

	environment, err := unmarshalStringList(o.Environment)
	if err != nil {
		return nil, err
	}

	importList, err := o.Import.Translate()
	if err != nil {
		return nil, err
	}

	exportList, err := o.Export.Translate()
	if err != nil {
		return nil, err
	}

	return &config.Override{
		Options:        options,
		Registries:     registries,
		Environment:    environment,
		ImportExcludes: importList.Excludes,
		ExportExcludes: exportList.Excludes,
	}, nil
}
