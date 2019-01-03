package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type RegistrySuite struct{}

func (s *RegistrySuite) TestTranslateRegistry(t sweet.T) {
	registry, err := translateRegistry(json.RawMessage(`{
		"type": "gcr",
		"key-file": "/etc/build-agent-key.json"
	}`))

	Expect(err).To(BeNil())
	Expect(registry.GetType()).To(Equal("gcr"))
	Expect(registry.(*config.GCRRegistry).KeyFile).To(Equal("/etc/build-agent-key.json"))
}

func (s *RegistrySuite) TestTranslateRegistryDefaultType(t sweet.T) {
	registry, err := translateRegistry(json.RawMessage(`{
		"server": "docker.io",
		"username": "admin",
		"password": "secret"
	}`))

	Expect(err).To(BeNil())
	Expect(registry.GetType()).To(Equal("server"))
	Expect(registry.(*config.ServerRegistry).Server).To(Equal("docker.io"))
	Expect(registry.(*config.ServerRegistry).Username).To(Equal("admin"))
	Expect(registry.(*config.ServerRegistry).Password).To(Equal("secret"))
}

func (s *RegistrySuite) TestTranslateRegistryUnknownType(t sweet.T) {
	_, err := translateRegistry(json.RawMessage(`{
		"type": "portus",
		"username": "admin"
	}`))

	Expect(err).To(MatchError("unknown registry type 'portus'"))
}

func (s *RegistrySuite) TestTranslateRegistryMalformedSchema(t sweet.T) {
	_, err := translateRegistry(json.RawMessage(`{
		"type": "ecr",
		"key-file": "/etc/build-agent-key.json"
	}`))

	Expect(err).To(MatchError("failed to validate registry: Additional property key-file is not allowed"))
}
