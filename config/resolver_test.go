package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ResolverSuite struct{}

func (s *ResolverSuite) TestExtend(t sweet.T) {
	config := &Config{
		Tasks: map[string]Task{
			"foo": &BuildTask{TaskMeta: TaskMeta{Name: "foo", Extends: ""}, Dockerfile: "foo"},
			"bar": &BuildTask{TaskMeta: TaskMeta{Name: "bar", Extends: "foo"}, Dockerfile: "bar"},
			"baz": &BuildTask{TaskMeta: TaskMeta{Name: "baz", Extends: "foo"}, Dockerfile: ""},
			"bnk": &BuildTask{TaskMeta: TaskMeta{Name: "bnk", Extends: "bar"}, Dockerfile: ""},
		},
	}

	Expect(config.Resolve()).To(BeNil())
	Expect(config.Tasks["foo"].(*BuildTask).Dockerfile).To(Equal("foo"))
	Expect(config.Tasks["bar"].(*BuildTask).Dockerfile).To(Equal("bar"))
	Expect(config.Tasks["baz"].(*BuildTask).Dockerfile).To(Equal("foo"))
	Expect(config.Tasks["bnk"].(*BuildTask).Dockerfile).To(Equal("bar"))
}

func (s *ResolverSuite) TestCycle(t sweet.T) {
	config := &Config{
		Tasks: map[string]Task{
			"foo": &RunTask{TaskMeta: TaskMeta{Name: "foo", Extends: "bar"}},
			"bar": &RunTask{TaskMeta: TaskMeta{Name: "bar", Extends: "foo"}},
		},
	}

	err := config.Resolve()
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("extension is cyclic"))
}
