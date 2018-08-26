package command

import (
	"fmt"

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

	builder := NewBuilder(builders, []string{"a", "b", "c"})
	args, err := builder.Build()

	Expect(err).To(BeNil())
	Expect(args).To(Equal([]string{"a", "b", "c", "x", "y", "z"}))
}

func (s *BuilderSuite) TestAddFlagValue(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { b.AddFlagValue("x", ""); return nil },
		func(b *Builder) error { b.AddFlagValue("y", "2"); return nil },
		func(b *Builder) error { b.AddFlagValue("z", "3"); return nil },
	}

	builder := NewBuilder(builders, []string{"a", "b", "c"})
	args, err := builder.Build()

	Expect(err).To(BeNil())
	Expect(args).To(Equal([]string{"a", "b", "c", "y", "2", "z", "3"}))
}

func (s *BuilderSuite) TestBuildFuncError(t sweet.T) {
	builders := []BuildFunc{
		func(b *Builder) error { return fmt.Errorf("utoh") },
	}

	builder := NewBuilder(builders, []string{"a", "b", "c"})
	_, err := builder.Build()
	Expect(err).To(MatchError("utoh"))
}
