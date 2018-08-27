package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type PushSuite struct{}

func (s *PushSuite) TestTranslate(t sweet.T) {
	task := &PushTask{
		Extends: "parent",
		Images:  json.RawMessage(`["i1", "i2", "i3"]`),
	}

	translated, err := task.Translate("push")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PushTask{
		TaskMeta: config.TaskMeta{
			Name:    "push",
			Extends: "parent",
		},
		Images: []string{"i1", "i2", "i3"},
	}))
}

func (s *PushSuite) TestTranslateStringLists(t sweet.T) {
	task := &PushTask{
		Extends: "parent",
		Images:  json.RawMessage(`"i1"`),
	}

	translated, err := task.Translate("push")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.PushTask{
		TaskMeta: config.TaskMeta{
			Name:    "push",
			Extends: "parent",
		},
		Images: []string{"i1"},
	}))
}
