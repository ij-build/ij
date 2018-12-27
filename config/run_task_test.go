package config

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RunTaskSuite struct{}

func (s *RunTaskSuite) TestExtend(t sweet.T) {
	parentHealthcheck := &Healthcheck{
		Command:     "parent-command",
		Interval:    time.Minute,
		Retries:     5,
		StartPeriod: time.Minute,
		Timeout:     time.Minute,
	}

	childHealthcheck := &Healthcheck{
		Command:     "child-command",
		Interval:    time.Second,
		Retries:     10,
		StartPeriod: time.Second,
		Timeout:     time.Second,
	}

	parent := &RunTask{
		TaskMeta: TaskMeta{
			Name:                "parent",
			Environment:         []string{"parent-env1"},
			RequiredEnvironment: []string{"parent-env2"},
		},
		Image:                  "parent-image",
		Command:                "parent-command",
		Shell:                  "parent-shell",
		Script:                 "parent-script",
		Entrypoint:             "parent-entrypoint",
		User:                   "parent-user",
		Workspace:              "parent-workspace",
		Hostname:               "parent-hostname",
		Detach:                 false,
		Healthcheck:            parentHealthcheck,
		ExportEnvironmentFiles: []string{"parent-exp1"},
	}

	child := &RunTask{
		TaskMeta: TaskMeta{
			Name:                "child",
			Extends:             "parent",
			Environment:         []string{"child-env1"},
			RequiredEnvironment: []string{"child-env2"},
		},
		Image:                  "child-image",
		Command:                "child-command",
		Shell:                  "child-shell",
		Script:                 "child-script",
		Entrypoint:             "child-entrypoint",
		User:                   "child-user",
		Workspace:              "child-workspace",
		Hostname:               "child-hostname",
		Detach:                 true,
		Healthcheck:            childHealthcheck,
		ExportEnvironmentFiles: []string{"child-exp1"},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Environment).To(Equal([]string{"parent-env1", "child-env1"}))
	Expect(child.RequiredEnvironment).To(Equal([]string{"parent-env2", "child-env2"}))
	Expect(child.Image).To(Equal("child-image"))
	Expect(child.Command).To(Equal("child-command"))
	Expect(child.Shell).To(Equal("child-shell"))
	Expect(child.Script).To(Equal("child-script"))
	Expect(child.Entrypoint).To(Equal("child-entrypoint"))
	Expect(child.User).To(Equal("child-user"))
	Expect(child.Workspace).To(Equal("child-workspace"))
	Expect(child.Hostname).To(Equal("child-hostname"))
	Expect(child.Detach).To(BeTrue())
	Expect(child.Healthcheck.Command).To(Equal("child-command"))
	Expect(child.Healthcheck.Interval).To(Equal(time.Second))
	Expect(child.Healthcheck.Retries).To(Equal(10))
	Expect(child.Healthcheck.StartPeriod).To(Equal(time.Second))
	Expect(child.Healthcheck.Timeout).To(Equal(time.Second))
	Expect(child.ExportEnvironmentFiles).To(ConsistOf("parent-exp1", "child-exp1"))
}

func (s *RunTaskSuite) TestExtendNoOverride(t sweet.T) {
	parentHealthcheck := &Healthcheck{
		Command:     "parent-command",
		Interval:    time.Minute,
		Retries:     5,
		StartPeriod: time.Minute,
		Timeout:     time.Minute,
	}

	parent := &RunTask{
		TaskMeta:    TaskMeta{Name: "parent"},
		Image:       "parent-image",
		Command:     "parent-command",
		Shell:       "parent-shell",
		Script:      "parent-script",
		Entrypoint:  "parent-entrypoint",
		User:        "parent-user",
		Workspace:   "parent-workspace",
		Hostname:    "parent-hostname",
		Detach:      true,
		Healthcheck: parentHealthcheck,
	}

	child := &RunTask{
		TaskMeta:    TaskMeta{Name: "child", Extends: "parent"},
		Healthcheck: &Healthcheck{},
	}

	Expect(child.Extend(parent)).To(BeNil())
	Expect(child.Image).To(Equal("parent-image"))
	Expect(child.Command).To(Equal("parent-command"))
	Expect(child.Shell).To(Equal("parent-shell"))
	Expect(child.Script).To(Equal("parent-script"))
	Expect(child.Entrypoint).To(Equal("parent-entrypoint"))
	Expect(child.User).To(Equal("parent-user"))
	Expect(child.Workspace).To(Equal("parent-workspace"))
	Expect(child.Hostname).To(Equal("parent-hostname"))
	Expect(child.Detach).To(BeTrue())
	Expect(child.Healthcheck.Command).To(Equal("parent-command"))
	Expect(child.Healthcheck.Interval).To(Equal(time.Minute))
	Expect(child.Healthcheck.Retries).To(Equal(5))
	Expect(child.Healthcheck.StartPeriod).To(Equal(time.Minute))
	Expect(child.Healthcheck.Timeout).To(Equal(time.Minute))
}

func (s *RunTaskSuite) TestExtendWrongType(t sweet.T) {
	parent := &BuildTask{TaskMeta: TaskMeta{Name: "parent"}}
	child := &RunTask{TaskMeta: TaskMeta{Name: "child", Extends: "parent"}}

	Expect(child.Extend(parent)).NotTo(BeNil())
}
