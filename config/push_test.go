package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PushSuite struct{}

func (s *PushSuite) TestExtend(t sweet.T) {
	parent := &PushTask{
		TaskMeta: TaskMeta{Name: "parent", Extends: ""},
		Images:   []string{"parent-i1"},
	}

	child := &PushTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
		Images:   []string{"child-i2", "child-i3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Images).To(ConsistOf("parent-i1", "child-i2", "child-i3"))
}

func (s *PushSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &PushTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
