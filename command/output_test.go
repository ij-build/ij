package command

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type OutputSuite struct{}

func (s *OutputSuite) TestLogProcessor(t sweet.T) {
	var (
		formats = make(chan string, 3)
		prefix  = logging.NewPrefix("a", "b", "c")
	)

	defer close(formats)

	logFunc := func(p *logging.Prefix, format string, _ ...interface{}) {
		Expect(p).To(Equal(prefix))
		formats <- format
	}

	processor := newLogProcessor(prefix, logFunc)
	processor.Process("foo")
	processor.Process("bar")
	processor.Process("baz")

	Expect(formats).To(Receive(Equal("foo")))
	Expect(formats).To(Receive(Equal("bar")))
	Expect(formats).To(Receive(Equal("baz")))
}

func (s *OutputSuite) TestStringProcessor(t sweet.T) {
	processor := newStringProcessor()
	processor.Process("foo")
	processor.Process("bar")
	processor.Process("baz")

	Expect(processor.String()).To(Equal("foo\nbar\nbaz\n"))
}
