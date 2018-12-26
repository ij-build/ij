package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RemoveTaskSuite struct{}

func (s *RemoveTaskSuite) TestExtend(t sweet.T) {
	parent := &RemoveTask{
		TaskMeta:    TaskMeta{Name: "parent", Extends: ""},
		Images:      []string{"parent-i1"},
		Environment: []string{"X=1", "Y=2"},
	}

	child := &RemoveTask{
		TaskMeta:    TaskMeta{Name: "child", Extends: "parent"},
		Images:      []string{"child-i2", "child-i3"},
		Environment: []string{"Y=3", "Z=4"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Images).To(ConsistOf("parent-i1", "child-i2", "child-i3"))
	Expect(child.Environment).To(Equal([]string{"X=1", "Y=2", "Y=3", "Z=4"}))
}

func (s *RemoveTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &RemoveTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}

func (s *RemoveTaskSuite) TestGetEnvironment(t sweet.T) {
	task := &RemoveTask{
		TaskMeta:    TaskMeta{Name: "task", Extends: ""},
		Environment: []string{"env1", "env2", "env3"},
	}

	Expect(task.GetEnvironment()).To(Equal([]string{"env1", "env2", "env3"}))
}
