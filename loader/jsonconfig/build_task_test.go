package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/ij-build/ij/config"
	. "github.com/onsi/gomega"
)

type BuildTaskSuite struct{}

func (s *BuildTaskSuite) TestTranslate(t sweet.T) {
	task := &BuildTask{
		Extends:             "parent",
		Environment:         json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
		RequiredEnvironment: []string{"X"},
		Dockerfile:          "dockerfile",
		Target:              "target",
		Tags:                json.RawMessage(`["t1", "t2", "t3"]`),
		Labels:              json.RawMessage(`["l1", "l2", "l3"]`),
	}

	translated, err := task.Translate("build")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.BuildTask{
		TaskMeta: config.TaskMeta{
			Name:                "build",
			Extends:             "parent",
			Environment:         []string{"X=1", "Y=2", "Z=3"},
			RequiredEnvironment: []string{"X"},
		},
		Dockerfile: "dockerfile",
		Target:     "target",
		Tags:       []string{"t1", "t2", "t3"},
		Labels:     []string{"l1", "l2", "l3"},
	}))
}

func (s *BuildTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &BuildTask{
		Extends:     "parent",
		Environment: json.RawMessage(`"X=1"`),
		Dockerfile:  "dockerfile",
		Target:      "target",
		Tags:        json.RawMessage(`"t1"`),
		Labels:      json.RawMessage(`"l1"`),
	}

	translated, err := task.Translate("build")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.BuildTask{
		TaskMeta: config.TaskMeta{
			Name:        "build",
			Extends:     "parent",
			Environment: []string{"X=1"},
		},
		Dockerfile: "dockerfile",
		Target:     "target",
		Tags:       []string{"t1"},
		Labels:     []string{"l1"},
	}))
}
