package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PushTaskSuite struct{}

func (s *PushTaskSuite) TestExtend(t sweet.T) {
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

func (s *PushTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &PushTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
