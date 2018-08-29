package config

type (
	Registry interface {
		GetType() string
	}

	GCRRegistry struct {
		KeyFile string
	}

	ECRRegistry struct {
		AccessKeyID     string
		SecretAccessKey string
		AccountID       string
		Region          string
		Role            string
	}

	ServerRegistry struct {
		Server       string
		Username     string
		Password     string
		PasswordFile string
	}
)

func (r *GCRRegistry) GetType() string    { return "gcr" }
func (r *ECRRegistry) GetType() string    { return "ecr" }
func (r *ServerRegistry) GetType() string { return "server" }
