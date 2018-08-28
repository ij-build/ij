package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type LogoutSuite struct{}

func (s *LogoutSuite) TestTranslate(t sweet.T) {
	task := &LogoutTask{
		Extends: "parent",
		Servers: json.RawMessage(`["s1", "s2", "s3"]`),
	}

	translated, err := task.Translate("logout")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.LogoutTask{
		TaskMeta: config.TaskMeta{
			Name:    "logout",
			Extends: "parent",
		},
		Servers: []string{"s1", "s2", "s3"},
	}))
}

func (s *LogoutSuite) TestTranslateStringLists(t sweet.T) {
	task := &LogoutTask{
		Extends: "parent",
		Servers: json.RawMessage(`"s1"`),
	}

	translated, err := task.Translate("logout")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.LogoutTask{
		TaskMeta: config.TaskMeta{
			Name:    "logout",
			Extends: "parent",
		},
		Servers: []string{"s1"},
	}))
}
