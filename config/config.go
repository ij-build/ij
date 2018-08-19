package config

type Config struct {
	Environment []string         `json:"environment"`
	Tasks       map[string]*Task `json:"tasks"`
	Plans       map[string]*Plan `json:"plans"`
}
