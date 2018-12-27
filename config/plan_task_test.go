package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PlanTaskSuite struct{}

func (s *PlanTaskSuite) TestExtend(t sweet.T) {
	parent := &PlanTask{
		TaskMeta: TaskMeta{
			Name:                "parent",
			Environment:         []string{"parent-env1"},
			RequiredEnvironment: []string{"parent-env2"},
		},
		Name: "parent-name",
	}

	child := &PlanTask{
		TaskMeta: TaskMeta{
			Name:                "child",
			Extends:             "parent",
			Environment:         []string{"child-env1"},
			RequiredEnvironment: []string{"child-env2"},
		},
		Name: "child-name",
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Environment).To(Equal([]string{"parent-env1", "child-env1"}))
	Expect(child.RequiredEnvironment).To(Equal([]string{"parent-env2", "child-env2"}))
	Expect(child.Name).To(Equal("child-name"))
}

func (s *PlanTaskSuite) TestExtendNoOverride(t sweet.T) {
	parent := &PlanTask{
		TaskMeta: TaskMeta{Name: "parent"},
		Name:     "parent-name",
	}

	child := &PlanTask{
		TaskMeta: TaskMeta{Name: "child", Extends: "parent"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Name).To(Equal("parent-name"))
}

func (s *PlanTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &BuildTask{TaskMeta: TaskMeta{Name: "parent"}}
	child := &PlanTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
