package jsonconfig

import (
	"encoding/json"

	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/util"
)

type Override struct {
	Options          *Options          `json:"options"`
	Registries       []json.RawMessage `json:"registries"`
	Environment      json.RawMessage   `json:"environment"`
	EnvironmentFiles json.RawMessage   `json:"env-file"`
	Import           *ImportFileList   `json:"import"`
	Export           *ExportFileList   `json:"export"`
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

	environment, err := util.UnmarshalStringList(o.Environment)
	if err != nil {
		return nil, err
	}

	environmentFiles, err := util.UnmarshalStringList(o.EnvironmentFiles)
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
		Options:          options,
		Registries:       registries,
		Environment:      environment,
		EnvironmentFiles: environmentFiles,
		ImportExcludes:   importList.Excludes,
		ExportExcludes:   exportList.Excludes,
		CleanExcludes:    exportList.CleanExcludes,
	}, nil
}
