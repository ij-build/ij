package jsonconfig

import (
	"encoding/json"

	"github.com/efritz/ij/config"
)

type Plan struct {
	Extend      bool            `json:"extend"`
	Stages      []*Stage        `json:"stages"`
	Environment json.RawMessage `json:"environment"`
}

func (p *Plan) Translate(name string) (*config.Plan, error) {
	stages := []*config.Stage{}
	for _, stage := range p.Stages {
		translated, err := stage.Translate()
		if err != nil {
			return nil, err
		}

		stages = append(stages, translated)
	}

	environment, err := unmarshalStringList(p.Environment)
	if err != nil {
		return nil, err
	}

	return &config.Plan{
		Name:        name,
		Extend:      p.Extend,
		Stages:      stages,
		Environment: environment,
	}, nil
}
