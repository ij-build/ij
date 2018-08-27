package command

//go:generate go-mockgen github.com/efritz/ij/logging -i Logger -o mock_logger_test.go -f

import (
	"context"

	"github.com/aphistic/sweet"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type RunnerSuite struct{}

func (s *RunnerSuite) TestRun(t sweet.T) {
	logger := NewMockLogger()

	err := newRunner(logger, true).Run(
		context.Background(),
		append(
			testArgs,
			"foo",
			"bar",
			"baz",
		),
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.InfoFuncCallCount()).To(Equal(3))
	Expect(logger.ErrorFuncCallCount()).To(Equal(0))

	params := logger.InfoFuncCallParams()
	Expect(params[0].Arg1).To(Equal("0 > foo"))
	Expect(params[1].Arg1).To(Equal("1 > bar"))
	Expect(params[2].Arg1).To(Equal("2 > baz"))
}

func (s *RunnerSuite) TestRunErrorOutput(t sweet.T) {
	logger := NewMockLogger()

	err := newRunner(logger, true).Run(
		context.Background(),
		append(
			testArgs,
			"foo",
			"FOO",
			"bar",
			"BAR",
			"baz",
			"BAZ",
		),
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.InfoFuncCallCount()).To(Equal(3))
	Expect(logger.ErrorFuncCallCount()).To(Equal(3))

	outParams := logger.InfoFuncCallParams()
	Expect(outParams[0].Arg1).To(Equal("0 > foo"))
	Expect(outParams[1].Arg1).To(Equal("2 > bar"))
	Expect(outParams[2].Arg1).To(Equal("4 > baz"))

	errParams := logger.ErrorFuncCallParams()
	Expect(errParams[0].Arg1).To(Equal("1 > FOO"))
	Expect(errParams[1].Arg1).To(Equal("3 > BAR"))
	Expect(errParams[2].Arg1).To(Equal("5 > BAZ"))
}

func (s *RunnerSuite) TestRunForOutput(t sweet.T) {
	runner := newRunner(logging.NilLogger, true)

	outText, errText, err := runner.RunForOutput(
		context.Background(),
		append(
			testArgs,
			"foo",
			"bar",
			"baz",
		),
	)

	Expect(err).To(BeNil())
	Expect(outText).To(Equal("0 > foo\n1 > bar\n2 > baz\n"))
	Expect(errText).To(BeEmpty())
}

func (s *RunnerSuite) TestRunForOutputErrorOutput(t sweet.T) {
	runner := newRunner(logging.NilLogger, true)

	outText, errText, err := runner.RunForOutput(context.Background(), append(
		testArgs,
		"foo",
		"FOO",
		"bar",
		"BAR",
		"baz",
		"BAZ",
	))

	Expect(err).To(BeNil())
	Expect(outText).To(Equal("0 > foo\n2 > bar\n4 > baz\n"))
	Expect(errText).To(Equal("1 > FOO\n3 > BAR\n5 > BAZ\n"))
}
