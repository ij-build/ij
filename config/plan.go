package config

import "fmt"

type Plan struct {
	Name        string   `json:"-"`
	Disabled    string   `json:"disabled,omitempty"`
	Extends     string   `json:"extends,omitempty"`
	Stages      []*Stage `json:"stages,omitempty"`
	Environment []string `json:"environment,omitempty"`
}

func (p *Plan) Merge(child *Plan) error {
	for _, stage := range child.Stages {
		if err := p.AddStage(stage); err != nil {
			return err
		}
	}

	p.Disabled = extendString(child.Disabled, p.Disabled)
	p.Environment = append(p.Environment, child.Environment...)
	return nil
}

func (p *Plan) AddStage(stage *Stage) error {
	if stage.BeforeStage != "" && stage.AfterStage != "" {
		return fmt.Errorf(
			"before_stage and after_stage declared in %s/%s",
			p.Name,
			stage.Name,
		)
	}

	if index := p.StageIndex(stage.Name); index >= 0 {
		if stage.BeforeStage != "" || stage.AfterStage != "" {
			return fmt.Errorf(
				"%s/%s exists in parent config, but before_stage or after_stage is also declared",
				p.Name,
				stage.Name,
			)
		}

		p.Stages[index] = stage
		return nil
	}

	targetStage := stage.Name

	if stage.BeforeStage != "" {
		if index := p.StageIndex(stage.BeforeStage); index >= 0 {
			p.InsertStage(stage, index)
			return nil
		}

		targetStage = stage.BeforeStage
	}

	if stage.AfterStage != "" {
		if index := p.StageIndex(stage.AfterStage); index >= 0 {
			p.InsertStage(stage, index+1)
			return nil
		}

		targetStage = stage.AfterStage
	}

	return fmt.Errorf(
		"stage %s/%s not declared in parent config",
		p.Name,
		targetStage,
	)
}

func (p *Plan) StageIndex(name string) int {
	if name == "" {
		return -1
	}

	for i, stage := range p.Stages {
		if stage.Name == name {
			return i
		}
	}

	return -1
}

func (p *Plan) InsertStage(stage *Stage, index int) {
	p.Stages = append(p.Stages, nil)
	copy(p.Stages[index+1:], p.Stages[index:])
	p.Stages[index] = stage
}
