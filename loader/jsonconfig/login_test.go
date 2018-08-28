package jsonconfig

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/ij/config"
	. "github.com/onsi/gomega"
)

type LoginSuite struct{}

func (s *LoginSuite) TestTranslate(t sweet.T) {
	task := &LoginTask{
		Extends:      "parent",
		Server:       "server",
		Username:     "username",
		Password:     "password",
		PasswordFile: "password-file",
	}

	translated, err := task.Translate("login")
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.LoginTask{
		TaskMeta: config.TaskMeta{
			Name:    "login",
			Extends: "parent",
		},
		Server:       "server",
		Username:     "username",
		Password:     "password",
		PasswordFile: "password-file",
	}))
}
