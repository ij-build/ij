package command

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type BuilderSuite struct{}

func (s *BuilderSuite) TestAddFlag(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { b.AddFlag("x"); return nil },
		func(b *Builder) error { b.AddFlag("y"); return nil },
		func(b *Builder) error { b.AddFlag("z"); return nil },
	}

	builder := NewBuilder([]string{"a", "b", "c"}, builders)
	args, stdin, err := builder.Build()

	Expect(err).To(BeNil())
	Expect(stdin).To(BeNil())
	Expect(args).To(Equal([]string{"a", "b", "c", "x", "y", "z"}))
}

func (s *BuilderSuite) TestAddFlagValue(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { b.AddFlagValue("x", ""); return nil },
		func(b *Builder) error { b.AddFlagValue("y", "2"); return nil },
		func(b *Builder) error { b.AddFlagValue("z", "3"); return nil },
	}

	builder := NewBuilder([]string{"a", "b", "c"}, builders)
	args, stdin, err := builder.Build()

	Expect(err).To(BeNil())
	Expect(stdin).To(BeNil())
	Expect(args).To(Equal([]string{"a", "b", "c", "y", "2", "z", "3"}))
}

func (s *BuilderSuite) TestAddArgs(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { b.AddArgs("q", "w"); return nil },
		func(b *Builder) error { b.AddFlagValue("y", "2"); return nil },
		func(b *Builder) error { b.AddFlagValue("z", "3"); return nil },
	}

	builder := NewBuilder([]string{"a", "b", "c"}, builders)
	args, stdin, err := builder.Build()

	Expect(err).To(BeNil())
	Expect(stdin).To(BeNil())
	Expect(args).To(Equal([]string{"a", "b", "c", "y", "2", "z", "3", "q", "w"}))
}

func (s *BuilderSuite) TestSetStdin(t sweet.T) {
	reader := ioutil.NopCloser(bytes.NewReader([]byte("ok")))

	builders := []BuildFunc{
		func(b *Builder) error { b.SetStdin(reader); return nil },
	}

	builder := NewBuilder([]string{"a", "b", "c"}, builders)
	_, stdin, err := builder.Build()
	Expect(err).To(BeNil())
	Expect(stdin).To(Equal(reader))
}

func (s *BuilderSuite) TestBuildFuncError(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { return fmt.Errorf("utoh") },
	}

	builder := NewBuilder([]string{"a", "b", "c"}, builders)
	_, _, err := builder.Build()
	Expect(err).To(MatchError("utoh"))
}
