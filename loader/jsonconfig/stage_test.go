package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/ij-build/ij/config"
	. "github.com/onsi/gomega"
)

type StageSuite struct{}

func (s *StageSuite) TestTranslate(t sweet.T) {
	stage := &Stage{
		Name:        "s",
		BeforeStage: "b",
		AfterStage:  "a",
		Tasks: []json.RawMessage{
			json.RawMessage(`"t1"`),
			json.RawMessage(`{"name": "t2"}`),
			json.RawMessage(`{"name": "t3", "environment": ["X=4", "Y=5"]}`),
		},
		Parallel:    true,
		Environment: json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
	}

	translated, err := stage.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Stage{
		Name:        "s",
		BeforeStage: "b",
		AfterStage:  "a",
		Tasks: []*config.StageTask{
			&config.StageTask{Name: "t1"},
			&config.StageTask{Name: "t2"},
			&config.StageTask{Name: "t3", Environment: []string{"X=4", "Y=5"}},
		},
		RunMode:     config.RunModeOnSuccess,
		Parallel:    true,
		Environment: []string{"X=1", "Y=2", "Z=3"},
	}))
}

func (s *StageSuite) TestTranslateStringLists(t sweet.T) {
	stage := &Stage{
		Environment: json.RawMessage(`"X=1"`),
	}

	translated, err := stage.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Stage{
		Tasks:       []*config.StageTask{},
		RunMode:     config.RunModeOnSuccess,
		Environment: []string{"X=1"},
	}))
}

func (s *StageSuite) TestTranslateNestedStringLists(t sweet.T) {
	stage := &Stage{
		Tasks: []json.RawMessage{
			json.RawMessage(`{"name": "t1", "environment": "X=1"}`),
		},
	}

	translated, err := stage.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.Stage{
		Tasks: []*config.StageTask{
			&config.StageTask{Name: "t1", Environment: []string{"X=1"}},
		},
		RunMode: config.RunModeOnSuccess,
	}))
}

func (s *StageSuite) TestTranslateRunMode(t sweet.T) {
	modes := map[string]config.RunMode{
		"on-success": config.RunModeOnSuccess,
		"on-failure": config.RunModeOnFailure,
		"always":     config.RunModeAlways,
	}

	for name, mode := range modes {
		stage := &Stage{
			Name:    "s",
			RunMode: name,
		}

		translated, err := stage.Translate()
		Expect(err).To(BeNil())
		Expect(translated).To(Equal(&config.Stage{
			Name:    "s",
			Tasks:   []*config.StageTask{},
			RunMode: mode,
		}))
	}
}
