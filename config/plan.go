package config

type Plan struct {
	Name        string
	Stages      []*Stage `json:"stages"`
	Environment []string `json:"environment"`
}

type Stage struct {
	Name        string        `json:"name"`
	Tasks       []*StageTasks `json:"tasks"`
	Concurrent  bool          `json:"concurrent"`
	Environment []string      `json:"environment"`
}

type StageTasks struct {
	Name        string   `json:"name"`
	Environment []string `json:"environment"`
}
