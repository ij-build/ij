package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PlanTaskSuite struct{}

func (s *PlanTaskSuite) TestExtend(t sweet.T) {
	parent := &PlanTask{
		TaskMeta:    TaskMeta{Name: "parent", Extends: ""},
		Name:        "parent-name",
		Environment: []string{"parent-env1"},
	}

	child := &PlanTask{
		TaskMeta:    TaskMeta{Name: "child", Extends: "parent"},
		Name:        "child-name",
		Environment: []string{"child-env1"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Name).To(Equal("child-name"))
	Expect(child.Environment).To(ConsistOf("parent-env1", "child-env1"))
}

func (s *PlanTaskSuite) TestExtendNoOverride(t sweet.T) {
	parent := &PlanTask{
		TaskMeta:    TaskMeta{Name: "parent", Extends: ""},
		Name:        "parent-name",
		Environment: []string{"parent-env1"},
	}

	child := &PlanTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Name).To(Equal("parent-name"))
	Expect(child.Environment).To(ConsistOf("parent-env1"))
}
func (s *PlanTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &BuildTask{TaskMeta: TaskMeta{Name: "parent", Extends: ""}}
	child := &PlanTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
