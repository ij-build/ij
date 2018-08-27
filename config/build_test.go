package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type BuildSuite struct{}

func (s *BuildSuite) TestExtend(t sweet.T) {
	parent := &BuildTask{
		TaskMeta:   TaskMeta{Name: "parent", Extends: ""},
		Dockerfile: "Dockerfile.parent",
		Tags:       []string{"parent-t1"},
		Labels:     []string{"parent-l1"},
		Arguments:  []string{"parent-a1"},
	}

	child := &BuildTask{
		TaskMeta:   TaskMeta{Name: "child", Extends: "parent"},
		Dockerfile: "Dockerfile.child",
		Tags:       []string{"child-t2", "child-t3"},
		Labels:     []string{"child-l2", "child-l3"},
		Arguments:  []string{"child-a2", "child-a3"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Dockerfile).To(Equal("Dockerfile.child"))
	Expect(child.Tags).To(ConsistOf("parent-t1", "child-t2", "child-t3"))
	Expect(child.Labels).To(ConsistOf("parent-l1", "child-l2", "child-l3"))
	Expect(child.Arguments).To(ConsistOf("parent-a1", "child-a2", "child-a3"))
}

func (s *BuildSuite) TestExtendNoOverwrite(t sweet.T) {
	parent := &BuildTask{
		TaskMeta:   TaskMeta{Name: "parent", Extends: ""},
		Dockerfile: "Dockerfile.parent",
	}

	child := &BuildTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Dockerfile).To(Equal("Dockerfile.parent"))
}

func (s *BuildSuite) TestExtendWrongType(t sweet.T) {
	parent := &RunTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &BuildTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
