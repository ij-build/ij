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
		Extends:       "parent",
		SSHIdentities: json.RawMessage(`"*"`),
		Environment:   []string{"X=1", "Y=2", "Z=3"},
		Imports:       json.RawMessage(`"."`),
		Exports:       json.RawMessage(`"**/junit*.xml"`),
		Excludes:      json.RawMessage(`"**/__pycache__"`),
		Workspace:     "/go/src/example.com",
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
		Extends:       "parent",
		SSHIdentities: []string{"*"},
		Environment:   []string{"X=1", "Y=2", "Z=3"},
		Imports:       []string{"."},
		Exports:       []string{"**/junit*.xml"},
		Excludes:      []string{"**/__pycache__"},
		Workspace:     "/go/src/example.com",
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

func (s *ConfigSuite) TestTranslateArrays(t sweet.T) {
	jsonConfig := &Config{
		SSHIdentities: json.RawMessage(`["fp1", "fp2"]`),
		Imports:       json.RawMessage(`["src", "test"]`),
		Exports:       json.RawMessage(`["*.txt", "*.go"]`),
		Excludes:      json.RawMessage(`["*.cache", "*.temp"]`),
		Tasks:         map[string]json.RawMessage{},
		Plans:         map[string]*Plan{},
		Metaplans:     map[string][]string{},
	}

	translated, err := jsonConfig.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Config{
		SSHIdentities: []string{"fp1", "fp2"},
		Imports:       []string{"src", "test"},
		Exports:       []string{"*.txt", "*.go"},
		Excludes:      []string{"*.cache", "*.temp"},
		Tasks:         map[string]config.Task{},
		Plans:         map[string]*config.Plan{},
		Metaplans:     map[string][]string{},
	}))
}
