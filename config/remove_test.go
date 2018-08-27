package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RemoveSuite struct{}

func (s *RemoveSuite) TestExtend(t sweet.T) {
	parent := &RemoveTask{
		TaskMeta: TaskMeta{Name: "parent", Extends: ""},
		Images:   []string{"parent-i1"},
	}

	child := &RemoveTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
		Images:   []string{"child-i2", "child-i3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Images).To(ConsistOf("parent-i1", "child-i2", "child-i3"))
}

func (s *RemoveSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &RemoveTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
