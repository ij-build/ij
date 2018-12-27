package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type BuildTaskSuite struct{}

func (s *BuildTaskSuite) TestExtend(t sweet.T) {
	parent := &BuildTask{
		TaskMeta: TaskMeta{
			Name:                "parent",
			Environment:         []string{"parent-env1"},
			RequiredEnvironment: []string{"parent-env2"},
		},
		Dockerfile: "Dockerfile.parent",
		Target:     "parent-target",
		Tags:       []string{"parent-t1"},
		Labels:     []string{"parent-l1"},
	}

	child := &BuildTask{
		TaskMeta: TaskMeta{
			Name:                "child",
			Extends:             "parent",
			Environment:         []string{"child-env1"},
			RequiredEnvironment: []string{"child-env2"},
		},
		Dockerfile: "Dockerfile.child",
		Target:     "child-target",
		Tags:       []string{"child-t2", "child-t3"},
		Labels:     []string{"child-l2", "child-l3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Environment).To(Equal([]string{"parent-env1", "child-env1"}))
	Expect(child.RequiredEnvironment).To(Equal([]string{"parent-env2", "child-env2"}))
	Expect(child.Dockerfile).To(Equal("Dockerfile.child"))
	Expect(child.Target).To(Equal("child-target"))
	Expect(child.Tags).To(ConsistOf("parent-t1", "child-t2", "child-t3"))
	Expect(child.Labels).To(ConsistOf("parent-l1", "child-l2", "child-l3"))

}

func (s *BuildTaskSuite) TestExtendNoOverwrite(t sweet.T) {
	parent := &BuildTask{
		TaskMeta:   TaskMeta{Name: "parent"},
		Dockerfile: "Dockerfile.parent",
		Target:     "parent-target",
	}

	child := &BuildTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Dockerfile).To(Equal("Dockerfile.parent"))
	Expect(child.Target).To(Equal("parent-target"))
}

func (s *BuildTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent"}}
	child := &BuildTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
