package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type PushTaskSuite struct{}

func (s *PushTaskSuite) TestTranslate(t sweet.T) {
	task := &PushTask{
		Extends:     "parent",
		Images:      json.RawMessage(`["i1", "i2", "i3"]`),
		Environment: json.RawMessage(`["X=1", "Y=2", "Z=3"]`),
	}

	translated, err := task.Translate("push")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PushTask{
		TaskMeta: config.TaskMeta{
			Name:    "push",
			Extends: "parent",
		},
		Images:      []string{"i1", "i2", "i3"},
		Environment: []string{"X=1", "Y=2", "Z=3"},
	}))
}

func (s *PushTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &PushTask{
		Extends:     "parent",
		Images:      json.RawMessage(`"i1"`),
		Environment: json.RawMessage(`"X=1"`),
	}

	translated, err := task.Translate("push")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PushTask{
		TaskMeta: config.TaskMeta{
			Name:    "push",
			Extends: "parent",
		},
		Images:      []string{"i1"},
		Environment: []string{"X=1"},
	}))
}
