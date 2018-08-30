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
		SSHIdentities:       json.RawMessage(`"*"`),
		HealthcheckInterval: Duration{time.Second * 10},
		Registries: []json.RawMessage{
			json.RawMessage(`{"server": "docker.io"}`),
			json.RawMessage(`{"type": "gcr", "key_file": "secret.key"}`),
		},
		Environment: json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		Import:      &FileList{Excludes: json.RawMessage(`"**/__pycache__"`)},
		Export:      &FileList{},
	}

	translated, err := jsonOverride.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Override{
		SSHIdentities:       []string{"*"},
		HealthcheckInterval: time.Second * 10,
		Registries: []config.Registry{
			&config.ServerRegistry{Server: "docker.io"},
			&config.GCRRegistry{KeyFile: "secret.key"},
		},
		Environment:    []string{"X=1", "Y=2", "Z=3"},
		ImportExcludes: []string{"**/__pycache__"},
	}))
}

func (s *OverrideSuite) TestTranslateStringLists(t sweet.T) {
	jsonOverride := &Override{
		Environment: json.RawMessage(`"X=1"`),
		Import:      &FileList{},
		Export:      &FileList{},
	}

	translated, err := jsonOverride.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Override{
		Registries:  []config.Registry{},
		Environment: []string{"X=1"},
	}))
}
