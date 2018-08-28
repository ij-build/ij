package jsonconfig

import (
	"github.com/efritz/ij/config"
)

type LoginTask struct {
	Extends      string `json:"extends"`
	Server       string `json:"server"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordFile string `json:"password_file"`
}

func (t *LoginTask) Translate(name string) (config.Task, error) {
	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.LoginTask{
		TaskMeta:     meta,
		Server:       t.Server,
		Username:     t.Username,
		Password:     t.Password,
		PasswordFile: t.PasswordFile,
	}, nil
}
