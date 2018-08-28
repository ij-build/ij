package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LogoutSuite struct{}

func (s *LogoutSuite) TestExtend(t sweet.T) {
	parent := &LogoutTask{
		TaskMeta: TaskMeta{Name: "parent", Extends: ""},
		Servers:  []string{"parent-i1"},
	}

	child := &LogoutTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
		Servers:  []string{"child-i2", "child-i3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Servers).To(ConsistOf("parent-i1", "child-i2", "child-i3"))
}

func (s *LogoutSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &LogoutTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
