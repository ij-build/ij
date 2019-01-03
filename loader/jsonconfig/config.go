package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type (
	Config struct {
		Extends          string                     `json:"extends"`
		Options          *Options                   `json:"options"`
		Registries       []json.RawMessage          `json:"registries"`
		Workspace        string                     `json:"workspace"`
		Environment      json.RawMessage            `json:"environment"`
		EnvironmentFiles json.RawMessage            `json:"env-file"`
		Import           *ImportFileList            `json:"import"`
		Export           *ExportFileList            `json:"export"`
		Tasks            map[string]json.RawMessage `json:"tasks"`
		Plans            map[string]*Plan           `json:"plans"`
		Metaplans        map[string][]string        `json:"metaplans"`
	}

	Options struct {
		SSHIdentities       json.RawMessage `json:"ssh-identities"`
		ForceSequential     bool            `json:"force-sequential"`
		HealthcheckInterval util.Duration   `json:"healthcheck-interval"`
	}

	ImportFileList struct {
		Files    json.RawMessage `json:"files"`
		Excludes json.RawMessage `json:"excludes"`
	}

	ExportFileList struct {
		Files         json.RawMessage `json:"files"`
		Excludes      json.RawMessage `json:"excludes"`
		CleanExcludes json.RawMessage `json:"clean-excludes"`
	}
)

func (c *Config) Translate(parent *config.Config) (*config.Config, error) {
	options, err := c.Options.Translate()
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

	environmentFiles, err := unmarshalStringList(c.EnvironmentFiles)
	if err != nil {
		return nil, err
	}

	importList, err := c.Import.Translate()
	if err != nil {
		return nil, err
	}

	exportList, err := c.Export.Translate()
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
		Extends:          c.Extends,
		Options:          options,
		Registries:       registries,
		Workspace:        c.Workspace,
		Environment:      environment,
		EnvironmentFiles: environmentFiles,
		Import:           importList,
		Export:           exportList,
		Tasks:            tasks,
		Plans:            plans,
		Metaplans:        c.Metaplans,
	}, nil
}

func (c *Options) Translate() (*config.Options, error) {
	sshIdentities, err := unmarshalStringList(c.SSHIdentities)
	if err != nil {
		return nil, err
	}

	return &config.Options{
		SSHIdentities:       sshIdentities,
		ForceSequential:     c.ForceSequential,
		HealthcheckInterval: c.HealthcheckInterval.Duration,
	}, nil
}

func (l *ImportFileList) Translate() (*config.ImportFileList, error) {
	files, err := unmarshalStringList(l.Files)
	if err != nil {
		return nil, err
	}

	excludes, err := unmarshalStringList(l.Excludes)
	if err != nil {
		return nil, err
	}

	return &config.ImportFileList{
		Files:    files,
		Excludes: excludes,
	}, nil
}

func (l *ExportFileList) Translate() (*config.ExportFileList, error) {
	files, err := unmarshalStringList(l.Files)
	if err != nil {
		return nil, err
	}

	excludes, err := unmarshalStringList(l.Excludes)
	if err != nil {
		return nil, err
	}

	cleanExcludes, err := unmarshalStringList(l.CleanExcludes)
	if err != nil {
		return nil, err
	}

	return &config.ExportFileList{
		Files:         files,
		Excludes:      excludes,
		CleanExcludes: cleanExcludes,
	}, nil
}
