package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type PlanTaskSuite struct{}

func (s *PlanTaskSuite) TestTranslate(t sweet.T) {
	task := &PlanTask{
		Extends:             "parent",
		Environment:         json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		RequiredEnvironment: []string{"X"},
		Name:                "rec",
	}

	translated, err := task.Translate("plan")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PlanTask{
		TaskMeta: config.TaskMeta{
			Name:                "plan",
			Extends:             "parent",
			Environment:         []string{"X=1", "Y=2", "Z=3"},
			RequiredEnvironment: []string{"X"},
		},
		Name: "rec",
	}))
}

func (s *PlanTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &PlanTask{
		Extends:     "parent",
		Environment: json.RawMessage(`"X=1"`),
		Name:        "rec",
	}

	translated, err := task.Translate("plan")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PlanTask{
		TaskMeta: config.TaskMeta{
			Name:        "plan",
			Extends:     "parent",
			Environment: []string{"X=1"},
		},
		Name: "rec",
	}))
}
