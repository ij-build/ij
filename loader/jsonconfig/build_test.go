package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type BuildSuite struct{}

func (s *BuildSuite) TestTranslate(t sweet.T) {
	task := &BuildTask{
		Extends:    "parent",
		Dockerfile: "dockerfile",
		Tags:       json.RawMessage(`["t1", "t2", "t3"]`),
		Labels:     json.RawMessage(`["l1", "l2", "l3"]`),
		Arguments:  json.RawMessage(`["a1", "a2", "a3"]`),
	}

	translated, err := task.Translate("build")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.BuildTask{
		TaskMeta: config.TaskMeta{
			Name:    "build",
			Extends: "parent",
		},
		Dockerfile: "dockerfile",
		Tags:       []string{"t1", "t2", "t3"},
		Labels:     []string{"l1", "l2", "l3"},
		Arguments:  []string{"a1", "a2", "a3"},
	}))
}

func (s *BuildSuite) TestTranslateStringLists(t sweet.T) {
	task := &BuildTask{
		Extends:    "parent",
		Dockerfile: "dockerfile",
		Tags:       json.RawMessage(`"t1"`),
		Labels:     json.RawMessage(`"l1"`),
		Arguments:  json.RawMessage(`"a1"`),
	}

	translated, err := task.Translate("build")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.BuildTask{
		TaskMeta: config.TaskMeta{
			Name:    "build",
			Extends: "parent",
		},
		Dockerfile: "dockerfile",
		Tags:       []string{"t1"},
		Labels:     []string{"l1"},
		Arguments:  []string{"a1"},
	}))
}
