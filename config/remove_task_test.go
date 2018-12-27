package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RemoveTaskSuite struct{}

func (s *RemoveTaskSuite) TestExtend(t sweet.T) {
	parent := &RemoveTask{
		TaskMeta: TaskMeta{
			Name:                "parent",
			Environment:         []string{"parent-env1"},
			RequiredEnvironment: []string{"parent-env2"},
		},
		Images: []string{"parent-i1"},
	}

	child := &RemoveTask{
		TaskMeta: TaskMeta{
			Name:                "child",
			Extends:             "parent",
			Environment:         []string{"child-env1"},
			RequiredEnvironment: []string{"child-env2"},
		},
		Images: []string{"child-i2", "child-i3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Environment).To(Equal([]string{"parent-env1", "child-env1"}))
	Expect(child.RequiredEnvironment).To(Equal([]string{"parent-env2", "child-env2"}))
	Expect(child.Images).To(ConsistOf("parent-i1", "child-i2", "child-i3"))
}

func (s *RemoveTaskSuite) TestExtendNoOverride(t sweet.T) {
	parent := &RemoveTask{
		TaskMeta: TaskMeta{Name: "parent"},
		Images:   []string{"i1", "i2"},
	}

	child := &RemoveTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Images).To(ConsistOf("i1", "i2"))
}

func (s *RemoveTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent"}}
	child := &RemoveTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
