package config

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestMerge(t sweet.T) {
	parent := &Config{
		Options: &Options{
			SSHIdentities: []string{"parent-ssh1"},
		},
		Registries:  []Registry{&ServerRegistry{Server: "parent.io"}},
		Environment: []string{"parent-env1"},
		Workspace:   "parent-workspace",
		Import: &FileList{
			Files:    []string{"parent-imp1"},
			Excludes: []string{"parent-exc1"},
		},
		Export: &FileList{
			Files: []string{"parent-exp1"},
		},
		Tasks: map[string]Task{
			"t1": &BuildTask{TaskMeta: TaskMeta{Name: "t1", Extends: ""}, Dockerfile: "a"},
			"t2": &BuildTask{TaskMeta: TaskMeta{Name: "t2", Extends: ""}, Dockerfile: "b"},
		},
		Plans: map[string]*Plan{
			"p1": &Plan{Name: "p1", Environment: []string{"X=1"}},
			"p2": &Plan{Name: "p2", Environment: []string{"X=2"}},
		},
		Metaplans: map[string][]string{
			"mp1": []string{"p1"},
			"mp2": []string{"p1", "p2"},
		},
	}

	child := &Config{
		Options: &Options{
			SSHIdentities:       []string{"child-ssh2", "child-ssh3"},
			ForceSequential:     true,
			HealthcheckInterval: time.Second * 10,
		},
		Registries:  []Registry{&ServerRegistry{Server: "child.io"}},
		Workspace:   "child-workspace",
		Environment: []string{"child-env2", "child-env3"},
		Import: &FileList{
			Files:    []string{"child-imp2", "child-imp3"},
			Excludes: []string{"child-exc2", "child-exc3"},
		},
		Export: &FileList{
			Files: []string{"child-exp2", "child-exp3"},
		},
		Tasks: map[string]Task{
			"t2": &BuildTask{TaskMeta: TaskMeta{Name: "t2", Extends: ""}, Dockerfile: "c"},
			"t3": &BuildTask{TaskMeta: TaskMeta{Name: "t3", Extends: ""}, Dockerfile: "d"},
		},
		Plans: map[string]*Plan{
			"p1": &Plan{Name: "p1", Extend: true, Environment: []string{"X=3"}},
			"p2": &Plan{Name: "p2", Environment: []string{"X=4"}},
			"p3": &Plan{Name: "p3", Environment: []string{"X=5"}},
		},
		Metaplans: map[string][]string{
			"mp2": []string{"p1", "p3"},
			"mp3": []string{"p2"},
		},
	}

	Expect(parent.Merge(child)).To(BeNil())
	Expect(parent.Options.SSHIdentities).To(ConsistOf("parent-ssh1", "child-ssh2", "child-ssh3"))
	Expect(parent.Options.ForceSequential).To(BeTrue())
	Expect(parent.Options.HealthcheckInterval).To(Equal(time.Second * 10))
	Expect(parent.Registries).To(ConsistOf(
		&ServerRegistry{Server: "parent.io"},
		&ServerRegistry{Server: "child.io"},
	))
	Expect(parent.Workspace).To(Equal("child-workspace"))
	Expect(parent.Environment).To(ConsistOf("parent-env1", "child-env2", "child-env3"))
	Expect(parent.Import.Files).To(ConsistOf("parent-imp1", "child-imp2", "child-imp3"))
	Expect(parent.Export.Files).To(ConsistOf("parent-exp1", "child-exp2", "child-exp3"))
	Expect(parent.Import.Excludes).To(ConsistOf("parent-exc1", "child-exc2", "child-exc3"))

	Expect(parent.Tasks).To(HaveLen(3))
	Expect(parent.Tasks["t1"].(*BuildTask).Dockerfile).To(Equal("a"))
	Expect(parent.Tasks["t2"].(*BuildTask).Dockerfile).To(Equal("c"))
	Expect(parent.Tasks["t3"].(*BuildTask).Dockerfile).To(Equal("d"))

	Expect(parent.Plans).To(HaveLen(3))
	Expect(parent.Plans["p1"].Environment).To(ConsistOf("X=1", "X=3"))
	Expect(parent.Plans["p2"].Environment).To(ConsistOf("X=4"))
	Expect(parent.Plans["p3"].Environment).To(ConsistOf("X=5"))

	Expect(parent.Metaplans).To(HaveLen(3))
	Expect(parent.Metaplans["mp1"]).To(Equal([]string{"p1"}))
	Expect(parent.Metaplans["mp2"]).To(Equal([]string{"p1", "p3"}))
	Expect(parent.Metaplans["mp3"]).To(Equal([]string{"p2"}))
}

func (s *ConfigSuite) TestMergeNoOverride(t sweet.T) {
	parent := &Config{
		Options:   &Options{},
		Workspace: "parent-workspace",
		Import:    &FileList{},
		Export:    &FileList{},
	}

	child := &Config{
		Options: &Options{},
		Import:  &FileList{},
		Export:  &FileList{},
	}

	Expect(parent.Merge(child)).To(BeNil())
	Expect(parent.Workspace).To(Equal("parent-workspace"))
}

func (s *ConfigSuite) TestApplyOverride(t sweet.T) {
	config := &Config{
		Options: &Options{
			SSHIdentities: []string{"config-ssh"},
		},
		Registries:  []Registry{&GCRRegistry{KeyFile: "config-gcr"}},
		Environment: []string{"X=1", "Y=2"},
		Import:      &FileList{Files: []string{"."}, Excludes: []string{".temp"}},
		Export:      &FileList{Files: []string{"*.py*"}, Excludes: nil},
	}

	override := &Override{
		Options: &Options{
			SSHIdentities:       []string{"override-ssh"},
			ForceSequential:     true,
			HealthcheckInterval: time.Second * 10,
		},
		Registries:     []Registry{&ECRRegistry{AccountID: "override-ecr"}},
		Environment:    []string{"X=3", "Z=2"},
		ImportExcludes: []string{"*.tmp"},
		ExportExcludes: []string{"*.pyc"},
	}

	config.ApplyOverride(override)
	Expect(config.Options.SSHIdentities).To(Equal([]string{"config-ssh", "override-ssh"}))
	Expect(config.Options.ForceSequential).To(BeTrue())
	Expect(config.Options.HealthcheckInterval).To(Equal(time.Second * 10))
	Expect(config.Registries).To(Equal([]Registry{
		&GCRRegistry{KeyFile: "config-gcr"},
		&ECRRegistry{AccountID: "override-ecr"},
	}))
	Expect(config.Environment).To(Equal([]string{"X=1", "Y=2", "X=3", "Z=2"}))
	Expect(config.Import.Excludes).To(Equal([]string{".temp", "*.tmp"}))
	Expect(config.Export.Excludes).To(Equal([]string{"*.pyc"}))
}

func (s *ConfigSuite) TestValidate(t sweet.T) {
	config := &Config{
		Tasks: map[string]Task{
			"t1": &RunTask{TaskMeta: TaskMeta{Name: "t1", Extends: ""}},
			"t2": &RunTask{TaskMeta: TaskMeta{Name: "t2", Extends: ""}},
		},

		Plans: map[string]*Plan{
			"p1": &Plan{
				Name: "p1",
				Stages: []*Stage{
					&Stage{
						Name: "s1",
						Tasks: []*StageTask{
							&StageTask{
								Name: "t1",
							},
							&StageTask{
								Name: "t2",
							},
						},
					},
				},
			},
			"p2": &Plan{
				Name: "p2",
				Stages: []*Stage{
					&Stage{
						Name: "s1",
						Tasks: []*StageTask{
							&StageTask{
								Name: "t2",
							},
						},
					},
				},
			},
		},

		Metaplans: map[string][]string{
			"m1": []string{"p1", "p2"},
		},
	}

	Expect(config.Validate()).To(BeNil())
}

func (s *ConfigSuite) TestValidateUnknownExtends(t sweet.T) {
	config := &Config{
		Tasks: map[string]Task{
			"foo": &BuildTask{TaskMeta: TaskMeta{Name: "foo", Extends: "unknown"}},
		},
	}

	Expect(config.Validate()).To(MatchError("unknown task name unknown referenced in task foo"))
}

func (s *ConfigSuite) TestValidateUnknownTask(t sweet.T) {
	config := &Config{
		Plans: map[string]*Plan{
			"foo": &Plan{
				Name: "foo",
				Stages: []*Stage{
					&Stage{
						Name: "bar",
						Tasks: []*StageTask{
							&StageTask{
								Name: "unknown",
							},
						},
					},
				},
			},
		},
	}

	Expect(config.Validate()).To(MatchError("unknown task name unknown referenced in foo/bar"))
}

func (s *ConfigSuite) TestValidatePlanDefinedTwice(t sweet.T) {
	config := &Config{
		Plans: map[string]*Plan{
			"dup": &Plan{},
		},
		Metaplans: map[string][]string{
			"dup": []string{"dup"},
		},
	}

	Expect(config.Validate()).To(MatchError("plan dup is defined twice"))
}
func (s *ConfigSuite) TestValidateUnknownPlan(t sweet.T) {
	config := &Config{
		Metaplans: map[string][]string{
			"default": []string{"unknown"},
		},
	}

	Expect(config.Validate()).To(MatchError("unknown plan name unknown referenced in metaplan default"))
}

func (s *ConfigSuite) TestIsPlanDefined(t sweet.T) {
	config := &Config{
		Plans: map[string]*Plan{
			"foo": &Plan{},
		},
		Metaplans: map[string][]string{
			"bar": []string{"foo"},
		},
	}

	Expect(config.IsPlanDefined("foo")).To(BeTrue())
	Expect(config.IsPlanDefined("bar")).To(BeTrue())
	Expect(config.IsPlanDefined("baz")).To(BeFalse())
}
