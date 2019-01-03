package jsonconfig

import (
	"encoding/json"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type OverrideSuite struct{}

func (s *OverrideSuite) TestTranslate(t sweet.T) {
	jsonOverride := &Override{
		Options: &Options{
			SSHIdentities:       json.RawMessage(`"*"`),
			HealthcheckInterval: Duration{time.Second * 10},
		},
		Registries: []json.RawMessage{
			json.RawMessage(`{"server": "docker.io"}`),
			json.RawMessage(`{"type": "gcr", "key-file": "secret.key"}`),
		},
		Environment: json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		Import:      &ImportFileList{Excludes: json.RawMessage(`"**/__pycache__"`)},
		Export:      &ExportFileList{CleanExcludes: json.RawMessage(`["*.txt", "*.pdf"]`)},
	}

	translated, err := jsonOverride.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Override{
		Options: &config.Options{
			SSHIdentities:       []string{"*"},
			HealthcheckInterval: time.Second * 10,
		},
		Registries: []config.Registry{
			&config.ServerRegistry{Server: "docker.io"},
			&config.GCRRegistry{
				Hostname: "gcr.io",
				KeyFile:  "secret.key",
			},
		},
		Environment:    []string{"X=1", "Y=2", "Z=3"},
		ImportExcludes: []string{"**/__pycache__"},
		CleanExcludes:  []string{"*.txt", "*.pdf"},
	}))
}

func (s *OverrideSuite) TestTranslateStringLists(t sweet.T) {
	jsonOverride := &Override{
		Options:     &Options{},
		Environment: json.RawMessage(`"X=1"`),
		Import:      &ImportFileList{},
		Export:      &ExportFileList{},
	}

	translated, err := jsonOverride.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Override{
		Options:     &config.Options{},
		Registries:  []config.Registry{},
		Environment: []string{"X=1"},
	}))
}
