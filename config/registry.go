package config

type (
	Registry interface {
		GetType() string
	}

	GCRRegistry struct {
		Hostname string `json:"hostname,omitempty"`
		Key      string `json:"key,omitempty"`
		KeyFile  string `json:"key-file,omitempty"`
	}

	ECRRegistry struct {
		AccessKeyID     string `json:"access-key-id,omitempty"`
		SecretAccessKey string `json:"secret-access-key,omitempty"`
		AccountID       string `json:"account-id,omitempty"`
		Region          string `json:"region,omitempty"`
		Role            string `json:"role,omitempty"`
	}

	ServerRegistry struct {
		Server       string `json:"server,omitempty"`
		Username     string `json:"username,omitempty"`
		Password     string `json:"password,omitempty"`
		PasswordFile string `json:"password-file,omitempty"`
	}
)

func (r *GCRRegistry) GetType() string    { return "gcr" }
func (r *ECRRegistry) GetType() string    { return "ecr" }
func (r *ServerRegistry) GetType() string { return "server" }
