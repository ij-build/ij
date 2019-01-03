package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestTranslate(t sweet.T) {
	jsonConfig := &Config{
		Extends: "parent",
		Registries: []json.RawMessage{
			json.RawMessage(`{"server": "docker.io"}`),
			json.RawMessage(`{"type": "gcr", "hostname": "eu.gcr.io", "key-file": "secret.key"}`),
		},
		Options: &Options{
			SSHIdentities: json.RawMessage(`"*"`),
		},
		Workspace:        "/go/src/example.com",
		Environment:      json.RawMessage(`"X=1"`),
		EnvironmentFiles: json.RawMessage(`"./env"`),
		Import: &ImportFileList{
			Files:    json.RawMessage(`"."`),
			Excludes: json.RawMessage(`"**/__pycache__"`),
		},
		Export: &ExportFileList{
			Files: json.RawMessage(`"**/junit*.xml"`),
		},
		Tasks: map[string]json.RawMessage{
			"t1": json.RawMessage(`{"image": "i1"}`),
			"t2": json.RawMessage(`{"image": "i2"}`),
		},
		Plans: map[string]*Plan{
			"p1": &Plan{Stages: []*Stage{&Stage{Tasks: []json.RawMessage{json.RawMessage(`"t1"`)}}}},
			"p2": &Plan{Stages: []*Stage{&Stage{Tasks: []json.RawMessage{json.RawMessage(`"t2"`)}}}},
		},
		Metaplans: map[string][]string{
			"default": []string{"a", "b"},
		},
	}

	translated, err := jsonConfig.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Config{
		Extends: "parent",
		Registries: []config.Registry{
			&config.ServerRegistry{Server: "docker.io"},
			&config.GCRRegistry{
				Hostname: "eu.gcr.io",
				KeyFile:  "secret.key",
			},
		},
		Options: &config.Options{
			SSHIdentities: []string{"*"},
		},
		Workspace:        "/go/src/example.com",
		Environment:      []string{"X=1"},
		EnvironmentFiles: []string{"./env"},
		Import: &config.ImportFileList{
			Files:    []string{"."},
			Excludes: []string{"**/__pycache__"},
		},
		Export: &config.ExportFileList{
			Files: []string{"**/junit*.xml"},
		},
		Tasks: map[string]config.Task{
			"t1": &config.RunTask{
				TaskMeta:    config.TaskMeta{Name: "t1"},
				Image:       "i1",
				Healthcheck: &config.Healthcheck{},
			},
			"t2": &config.RunTask{
				TaskMeta:    config.TaskMeta{Name: "t2"},
				Image:       "i2",
				Healthcheck: &config.Healthcheck{},
			},
		},
		Plans: map[string]*config.Plan{
			"p1": &config.Plan{
				Name: "p1",
				Stages: []*config.Stage{
					&config.Stage{
						Tasks:   []*config.StageTask{&config.StageTask{Name: "t1"}},
						RunMode: config.RunModeOnSuccess,
					},
				},
			},
			"p2": &config.Plan{
				Name: "p2",
				Stages: []*config.Stage{
					&config.Stage{
						Tasks:   []*config.StageTask{&config.StageTask{Name: "t2"}},
						RunMode: config.RunModeOnSuccess,
					},
				},
			},
		},
		Metaplans: map[string][]string{
			"default": []string{"a", "b"},
		},
	}))
}

func (s *ConfigSuite) TestTranslateStringLists(t sweet.T) {
	jsonConfig := &Config{
		Options: &Options{
			SSHIdentities: json.RawMessage(`["fp1", "fp2"]`),
		},
		Environment:      json.RawMessage(`["X=1", "Y=2"]`),
		EnvironmentFiles: json.RawMessage(`["./env", "./env.secret"]`),
		Import: &ImportFileList{
			Files:    json.RawMessage(`["src", "test"]`),
			Excludes: json.RawMessage(`["*.cache", "*.temp"]`),
		},
		Export: &ExportFileList{
			Files: json.RawMessage(`["*.txt", "*.go"]`),
		},
		Tasks:     map[string]json.RawMessage{},
		Plans:     map[string]*Plan{},
		Metaplans: map[string][]string{},
	}

	translated, err := jsonConfig.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Config{
		Options: &config.Options{
			SSHIdentities: []string{"fp1", "fp2"},
		},
		Registries:       []config.Registry{},
		Environment:      []string{"X=1", "Y=2"},
		EnvironmentFiles: []string{"./env", "./env.secret"},
		Import: &config.ImportFileList{
			Files:    []string{"src", "test"},
			Excludes: []string{"*.cache", "*.temp"},
		},
		Export: &config.ExportFileList{
			Files: []string{"*.txt", "*.go"},
		},
		Tasks:     map[string]config.Task{},
		Plans:     map[string]*config.Plan{},
		Metaplans: map[string][]string{},
	}))
}
