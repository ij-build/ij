package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LoginSuite struct{}

func (s *LoginSuite) TestExtend(t sweet.T) {
	parent := &LoginTask{
		TaskMeta:     TaskMeta{Name: "parent", Extends: ""},
		Server:       "parent-server",
		Username:     "parent-username",
		Password:     "parent-password",
		PasswordFile: "parent-password-file",
	}

	child := &LoginTask{
		TaskMeta:     TaskMeta{Name: "child", Extends: "parent"},
		Server:       "child-server",
		Username:     "child-username",
		Password:     "child-password",
		PasswordFile: "child-password-file",
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Server).To(Equal("child-server"))
	Expect(child.Username).To(Equal("child-username"))
	Expect(child.Password).To(Equal("child-password"))
	Expect(child.PasswordFile).To(Equal("child-password-file"))
}

func (s *LoginSuite) TestExtendNoOverride(t sweet.T) {
	parent := &LoginTask{
		TaskMeta:     TaskMeta{Name: "parent", Extends: ""},
		Server:       "parent-server",
		Username:     "parent-username",
		Password:     "parent-password",
		PasswordFile: "parent-password-file",
	}

	child := &LoginTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Server).To(Equal("parent-server"))
	Expect(child.Username).To(Equal("parent-username"))
	Expect(child.Password).To(Equal("parent-password"))
	Expect(child.PasswordFile).To(Equal("parent-password-file"))
}

func (s *LoginSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &LoginTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
