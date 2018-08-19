package config

type Task struct {
	Name                string
	Image               string   `json:"image"`
	Environment         []string `json:"environment"`
	RequiredEnvironment []string `json:"required_environment"`
	Command             string   `json:"command"`
	Script              string   `json:"script"`
	Shell               string   `json:"shell"`
	// TODO - Hostname
	// TODO - Detach
	// TODO - HealthCheck
}
