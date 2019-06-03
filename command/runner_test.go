package command

//go:generate go-mockgen -f github.com/efritz/ij/logging -i Logger -o mock_logger_test.go

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	"github.com/efritz/ij/logging"
	. "github.com/onsi/gomega"
)

type RunnerSuite struct{}

func (s *RunnerSuite) TestRun(t sweet.T) {
	logger := NewMockLogger()

	args := append(
		testArgs,
		"foo",
		"bar",
		"baz",
	)

	err := newRunner(logger, true).Run(
		context.Background(),
		args,
		nil,
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.InfoFunc).To(BeCalledN(3))
	Expect(logger.ErrorFunc).NotTo(BeCalled())
	Expect(logger.InfoFunc.History()[0].Arg2[0]).To(Equal("0 > foo"))
	Expect(logger.InfoFunc.History()[1].Arg2[0]).To(Equal("1 > bar"))
	Expect(logger.InfoFunc.History()[2].Arg2[0]).To(Equal("2 > baz"))
}

func (s *RunnerSuite) TestRunWithStdin(t sweet.T) {
	logger := NewMockLogger()

	args := append(
		testArgs,
		"foo",
		"bar",
		"baz",
	)

	err := newRunner(logger, true).Run(
		context.Background(),
		args,
		ioutil.NopCloser(bytes.NewReader([]byte("XXX"))),
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.InfoFunc).To(BeCalledN(4))
	Expect(logger.ErrorFunc).NotTo(BeCalled())
	Expect(logger.InfoFunc.History()[0].Arg2[0]).To(Equal("x > XXX"))
	Expect(logger.InfoFunc.History()[1].Arg2[0]).To(Equal("0 > foo"))
	Expect(logger.InfoFunc.History()[2].Arg2[0]).To(Equal("1 > bar"))
	Expect(logger.InfoFunc.History()[3].Arg2[0]).To(Equal("2 > baz"))
}

func (s *RunnerSuite) TestRunWithMaskedSecrets(t sweet.T) {
	logger := NewMockLogger()

	args := append(
		testArgs,
		"arg=plaintext",
		"api_password=secret",
		"AWS_SECRET_PASSWORD=hidden",
	)

	expectedArgs := append(
		testArgs,
		"arg=plaintext",
		"api_password=*****",
		"AWS_SECRET_PASSWORD=*****",
	)

	err := newRunner(logger, true).Run(
		context.Background(),
		args,
		nil,
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.DebugFunc).To(BeCalledOnce())
	Expect(logger.DebugFunc.History()[0].Arg2).To(Equal([]interface{}{strings.Join(expectedArgs, " ")}))
}

func (s *RunnerSuite) TestRunErrorOutput(t sweet.T) {
	logger := NewMockLogger()

	args := append(
		testArgs,
		"foo",
		"FOO",
		"bar",
		"BAR",
		"baz",
		"BAZ",
	)

	err := newRunner(logger, true).Run(
		context.Background(),
		args,
		nil,
		nil,
	)

	Expect(err).To(BeNil())
	Expect(logger.InfoFunc).To(BeCalledN(3))
	Expect(logger.ErrorFunc).To(BeCalledN(3))
	Expect(logger.InfoFunc.History()[0].Arg2[0]).To(Equal("0 > foo"))
	Expect(logger.InfoFunc.History()[1].Arg2[0]).To(Equal("2 > bar"))
	Expect(logger.InfoFunc.History()[2].Arg2[0]).To(Equal("4 > baz"))
	Expect(logger.ErrorFunc.History()[0].Arg2[0]).To(Equal("1 > FOO"))
	Expect(logger.ErrorFunc.History()[1].Arg2[0]).To(Equal("3 > BAR"))
	Expect(logger.ErrorFunc.History()[2].Arg2[0]).To(Equal("5 > BAZ"))
}

func (s *RunnerSuite) TestRunForOutput(t sweet.T) {
	runner := newRunner(logging.NilLogger, true)

	args := append(
		testArgs,
		"foo",
		"bar",
		"baz",
	)

	outText, errText, err := runner.RunForOutput(
		context.Background(),
		args,
		nil,
	)

	Expect(err).To(BeNil())
	Expect(outText).To(Equal("0 > foo\n1 > bar\n2 > baz\n"))
	Expect(errText).To(BeEmpty())
}

func (s *RunnerSuite) TestRunForOutputErrorOutput(t sweet.T) {
	runner := newRunner(logging.NilLogger, true)

	args := append(
		testArgs,
		"foo",
		"FOO",
		"bar",
		"BAR",
		"baz",
		"BAZ",
	)

	outText, errText, err := runner.RunForOutput(
		context.Background(),
		args,
		nil,
	)

	Expect(err).To(BeNil())
	Expect(outText).To(Equal("0 > foo\n2 > bar\n4 > baz\n"))
	Expect(errText).To(Equal("1 > FOO\n3 > BAR\n5 > BAZ\n"))
}
