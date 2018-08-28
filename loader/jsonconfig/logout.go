package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type LogoutTask struct {
	Extends string          `json:"extends"`
	Servers json.RawMessage `json:"servers"`
}

func (t *LogoutTask) Translate(name string) (config.Task, error) {
	servers, err := unmarshalStringList(t.Servers)
	if err != nil {
		return nil, err
	}

	meta := config.TaskMeta{
		Name:    name,
		Extends: t.Extends,
	}

	return &config.LogoutTask{
		TaskMeta: meta,
		Servers:  servers,
	}, nil
}
