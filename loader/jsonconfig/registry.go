package jsonconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ij-build/ij/config"
	"github.com/ij-build/ij/loader/schema"
)

type (
	RegistryTypeHint struct {
		Type string `json:"type"`
	}

	Registry interface {
		Translate() config.Registry
	}

	ServerRegistry struct {
		Server       string `json:"server"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		PasswordFile string `json:"password-file"`
	}

	GCRRegistry struct {
		Hostname string `json:"hostname"`
		Key      string `json:"key"`
		KeyFile  string `json:"key-file"`
	}

	ECRRegistry struct {
		AccessKeyID     string `json:"access-key-id"`
		SecretAccessKey string `json:"secret-access-key"`
		AccountID       string `json:"account-id"`
		Region          string `json:"region"`
		Role            string `json:"role"`
	}
)

const (
	DefaultECRRegion   = "us-east-1"
	DefaultGCRHostname = "gcr.io"
)

func translateRegistry(data json.RawMessage) (config.Registry, error) {
	typeHint := &RegistryTypeHint{Type: "server"}

	if err := json.Unmarshal(data, typeHint); err != nil {
		return nil, err
	}

	structMap := map[string]Registry{
		"gcr":    &GCRRegistry{},
		"ecr":    &ECRRegistry{},
		"server": &ServerRegistry{},
	}

	registry, ok := structMap[typeHint.Type]
	if !ok {
		return nil, fmt.Errorf("unknown registry type '%s'", typeHint.Type)
	}

	assetName := fmt.Sprintf(
		"schema/registry-%s.yaml",
		strings.Replace(typeHint.Type, "-", "_", -1),
	)

	if err := schema.Validate(assetName, data); err != nil {
		return nil, fmt.Errorf("failed to validate registry: %s", err.Error())
	}

	if err := json.Unmarshal(data, registry); err != nil {
		return nil, err
	}

	return registry.Translate(), nil
}

func (r *ServerRegistry) Translate() config.Registry {
	return &config.ServerRegistry{
		Server:       r.Server,
		Username:     r.Username,
		Password:     r.Password,
		PasswordFile: r.PasswordFile,
	}
}

func (r *GCRRegistry) Translate() config.Registry {
	if r.Hostname == "" {
		r.Hostname = DefaultGCRHostname
	}

	return &config.GCRRegistry{
		Hostname: r.Hostname,
		Key:      r.Key,
		KeyFile:  r.KeyFile,
	}
}

func (r *ECRRegistry) Translate() config.Registry {
	if r.Region == "" {
		r.Region = DefaultECRRegion
	}

	return &config.ECRRegistry{
		AccessKeyID:     r.AccessKeyID,
		SecretAccessKey: r.SecretAccessKey,
		AccountID:       r.AccountID,
		Region:          r.Region,
		Role:            r.Role,
	}
}
