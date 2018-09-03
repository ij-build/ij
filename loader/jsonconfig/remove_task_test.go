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
		Extends: "parent",
		Images:  json.RawMessage(`["i1", "i2", "i3"]`),
	}

	translated, err := task.Translate("remove")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RemoveTask{
		TaskMeta: config.TaskMeta{
			Name:    "remove",
			Extends: "parent",
		},
		Images: []string{"i1", "i2", "i3"},
	}))
}

func (s *RemoveTaskSuite) TestTranslateStringLists(t sweet.T) {
	task := &RemoveTask{
		Extends: "parent",
		Images:  json.RawMessage(`"i1"`),
	}

	translated, err := task.Translate("remove")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.RemoveTask{
		TaskMeta: config.TaskMeta{
			Name:    "remove",
			Extends: "parent",
		},
		Images: []string{"i1"},
	}))
}
