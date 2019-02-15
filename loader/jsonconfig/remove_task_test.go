package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type RemoveTaskSuite struct{}

func (s *RemoveTaskSuite) TestTranslate(t sweet.T) {
	task := &RemoveTask{
		Extends:             "parent",
		Environment:         json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		RequiredEnvironment: []string{"X"},
		Images:              json.RawMessage(`["i1", "i2", "i3"]`),
		IncludeBuilt:        true,
	}

	translated, err := task.Translate("remove")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RemoveTask{
		TaskMeta: config.TaskMeta{
			Name:                "remove",
			Extends:             "parent",
			Environment:         []string{"X=1", "Y=2", "Z=3"},
			RequiredEnvironment: []string{"X"},
		},
		Images:       []string{"i1", "i2", "i3"},
		IncludeBuilt: true,
	}))
}

func (s *RemoveTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &RemoveTask{
		Extends:      "parent",
		Environment:  json.RawMessage(`"X=1"`),
		Images:       json.RawMessage(`"i1"`),
		IncludeBuilt: true,
	}

	translated, err := task.Translate("remove")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RemoveTask{
		TaskMeta: config.TaskMeta{
			Name:        "remove",
			Extends:     "parent",
			Environment: []string{"X=1"},
		},
		Images:       []string{"i1"},
		IncludeBuilt: true,
	}))
}
