package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type OverrideSuite struct{}

func (s *OverrideSuite) TestTranslate(t sweet.T) {
	jsonOverride := &Override{
		Registries: []json.RawMessage{
			json.RawMessage(`{"server": "docker.io"}`),
			json.RawMessage(`{"type": "gcr", "key_file": "secret.key"}`),
		},
		SSHIdentities: json.RawMessage(`"*"`),
		Environment:   []string{"X=1", "Y=2", "Z=3"},
		Excludes:      json.RawMessage(`"**/__pycache__"`),
	}

	translated, err := jsonOverride.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Override{
		Registries: []config.Registry{
			&config.ServerRegistry{Server: "docker.io"},
			&config.GCRRegistry{KeyFile: "secret.key"},
		},
		SSHIdentities: []string{"*"},
		Environment:   []string{"X=1", "Y=2", "Z=3"},
		Excludes:      []string{"**/__pycache__"},
	}))
}
