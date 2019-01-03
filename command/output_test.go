package command

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type OutputSuite struct{}

func (s *OutputSuite) TestLogProcessor(t sweet.T) {
	var (
		lines  = make(chan interface{}, 3)
		prefix = logging.NewPrefix("a", "b", "c")
	)

	defer close(lines)

	logFunc := func(p *logging.Prefix, format string, args ...interface{}) {
		Expect(p).To(Equal(prefix))
		lines <- args[0]
	}

	processor := newLogProcessor(prefix, logFunc)
	processor.Process("foo")
	processor.Process("bar")
	processor.Process("baz")

	Expect(lines).To(Receive(Equal("foo")))
	Expect(lines).To(Receive(Equal("bar")))
	Expect(lines).To(Receive(Equal("baz")))
}

func (s *OutputSuite) TestStringProcessor(t sweet.T) {
	processor := newStringProcessor()
	processor.Process("foo")
	processor.Process("bar")
	processor.Process("baz")

	Expect(processor.String()).To(Equal("foo\nbar\nbaz\n"))
}
