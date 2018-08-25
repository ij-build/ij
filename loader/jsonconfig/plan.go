package jsonconfig

import (
	"github.com/efritz/ij/config"
)

type Plan struct {
	Extend      bool     `json:"extend"`
	Stages      []*Stage `json:"stages"`
	Environment []string `json:"environment"`
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

	return &config.Plan{
		Name:        name,
		Extend:      p.Extend,
		Stages:      stages,
		Environment: p.Environment,
	}, nil
}
