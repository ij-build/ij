package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PlanSuite struct{}

func (s *PlanSuite) TestMerge(t sweet.T) {
	parent := &Plan{
		Disabled: "${PARENT_DISABLED}",
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "b"},
			&Stage{Name: "c"},
		},
		Environment: []string{
			"X=1",
			"Y=2",
		},
	}

	child := &Plan{
		Disabled: "${CHILD_DISABLED}",
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "d", BeforeStage: "a"},
			&Stage{Name: "e", AfterStage: "c"},
		},
		Environment: []string{
			"X=4",
			"Z=3",
		},
	}

	Expect(parent.Merge(child)).To(BeNil())
	Expect(parent.Disabled).To(Equal("${CHILD_DISABLED}"))
	Expect(parent.Stages).To(HaveLen(5))
	Expect(parent.Stages[0].Name).To(Equal("d"))
	Expect(parent.Stages[1].Name).To(Equal("a"))
	Expect(parent.Stages[2].Name).To(Equal("b"))
	Expect(parent.Stages[3].Name).To(Equal("c"))
	Expect(parent.Stages[4].Name).To(Equal("e"))
	Expect(parent.Environment).To(Equal([]string{"X=1", "Y=2", "X=4", "Z=3"}))
}

func (s *PlanSuite) TestAddStageOverwrite(t sweet.T) {
	plan := &Plan{
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "b"},
			&Stage{Name: "c"},
		},
	}

	stage := &Stage{
		Name:        "b",
		Environment: []string{"overwritten"},
	}

	Expect(plan.AddStage(stage)).To(BeNil())
	Expect(plan.Stages).To(HaveLen(3))
	Expect(plan.Stages[0].Name).To(Equal("a"))
	Expect(plan.Stages[1].Name).To(Equal("b"))
	Expect(plan.Stages[2].Name).To(Equal("c"))
	Expect(plan.Stages[1].Environment).To(ConsistOf("overwritten"))
}

func (s *PlanSuite) TestAddStageBefore(t sweet.T) {
	plan := &Plan{
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "b"},
			&Stage{Name: "c"},
		},
	}

	stage := &Stage{
		Name:        "d",
		BeforeStage: "b",
	}

	Expect(plan.AddStage(stage)).To(BeNil())
	Expect(plan.Stages).To(HaveLen(4))
	Expect(plan.Stages[0].Name).To(Equal("a"))
	Expect(plan.Stages[1].Name).To(Equal("d"))
	Expect(plan.Stages[2].Name).To(Equal("b"))
	Expect(plan.Stages[3].Name).To(Equal("c"))
}

func (s *PlanSuite) TestAddStageAfter(t sweet.T) {
	plan := &Plan{
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "b"},
			&Stage{Name: "c"},
		},
	}

	stage := &Stage{
		Name:       "d",
		AfterStage: "b",
	}

	Expect(plan.AddStage(stage)).To(BeNil())
	Expect(plan.Stages).To(HaveLen(4))
	Expect(plan.Stages[0].Name).To(Equal("a"))
	Expect(plan.Stages[1].Name).To(Equal("b"))
	Expect(plan.Stages[2].Name).To(Equal("d"))
	Expect(plan.Stages[3].Name).To(Equal("c"))
}

func (s *PlanSuite) TestAddStageAmbiguous(t sweet.T) {
	plan := &Plan{
		Name: "p",
	}

	err := plan.AddStage(&Stage{
		Name:        "s",
		BeforeStage: "b",
		AfterStage:  "a",
	})

	Expect(err).To(MatchError("before_stage and after_stage declared in p/s"))
}

func (s *PlanSuite) TestAddStageAmbiguousOverwrite(t sweet.T) {
	plan := &Plan{
		Name: "p",
		Stages: []*Stage{
			&Stage{Name: "s"},
		},
	}

	err := plan.AddStage(&Stage{
		Name:        "s",
		BeforeStage: "s",
	})

	Expect(err).To(MatchError("p/s exists in parent config, but before_stage or after_stage is also declared"))
}

func (s *PlanSuite) TestAddStageMissing(t sweet.T) {
	plan := &Plan{
		Name: "p",
	}

	err := plan.AddStage(&Stage{
		Name:        "s",
		BeforeStage: "b",
	})

	Expect(err).To(MatchError("stage p/b not declared in parent config"))
}

func (s *PlanSuite) TestStageIndex(t sweet.T) {
	plan := &Plan{
		Stages: []*Stage{
			&Stage{Name: "a"},
			&Stage{Name: "b"},
			&Stage{Name: "c"},
		},
	}

	Expect(plan.StageIndex("a")).To(Equal(0))
	Expect(plan.StageIndex("b")).To(Equal(1))
	Expect(plan.StageIndex("c")).To(Equal(2))
	Expect(plan.StageIndex("d")).To(Equal(-1))
}

func (s *PlanSuite) TestInsertStage(t sweet.T) {
	plan := &Plan{}
	plan.InsertStage(&Stage{Name: "a"}, 0) // a
	plan.InsertStage(&Stage{Name: "b"}, 0) // b a
	plan.InsertStage(&Stage{Name: "c"}, 0) // c b a
	plan.InsertStage(&Stage{Name: "d"}, 1) // c d b a
	plan.InsertStage(&Stage{Name: "e"}, 4) // c d b a e

	Expect(plan.Stages).To(HaveLen(5))
	Expect(plan.Stages[0].Name).To(Equal("c"))
	Expect(plan.Stages[1].Name).To(Equal("d"))
	Expect(plan.Stages[2].Name).To(Equal("b"))
	Expect(plan.Stages[3].Name).To(Equal("a"))
	Expect(plan.Stages[4].Name).To(Equal("e"))
}
