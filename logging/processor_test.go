package logging

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type ProcessorSuite struct {
	oldStdout *os.File
	oldStderr *os.File
}

func (s *ProcessorSuite) SetUpSuite() {
	outFile, _ := ioutil.TempFile("", "ij-stdout")
	s.oldStdout = os.Stdout
	os.Stdout = outFile

	errFile, _ := ioutil.TempFile("", "ij-stderr")
	s.oldStderr = os.Stderr
	os.Stderr = errFile
}

func (s *ProcessorSuite) TearDownSuite() {
	os.Stdout.Close()
	os.Stdout = s.oldStdout
	s.oldStdout = nil

	os.Stderr.Close()
	os.Stderr = s.oldStderr
	s.oldStderr = nil
}

func (s *ProcessorSuite) TestBasic(t sweet.T) {
	clock := glock.NewMockClock()
	processor := newProcessor(true, false, clock)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outfile = &logWriter{}
		errfile = &logWriter{}
		logger  = processor.Logger(outfile, errfile, true)
	)

	logger.Debug(nil, "> %s", "a")
	logger.Info(NewPrefix("x", "y", "z"), "> %s", "b")
	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "c")
	processor.Shutdown()

	lines := []string{
		"2018-08-26 09:29:03.123 | > a\n",
		"2018-08-26 09:29:03.123 | x/y/z: > b\n",
		"2018-08-26 09:29:03.123 | x/y/z: > c\n",
	}

	Expect(outfile.String()).To(Equal(strings.Join(lines, "")))
	Expect(errfile.String()).To(BeEmpty())
}

func (s *ProcessorSuite) TestErrors(t sweet.T) {
	clock := glock.NewMockClock()
	processor := newProcessor(true, false, clock)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outfile = &logWriter{}
		errfile = &logWriter{}
		logger  = processor.Logger(outfile, errfile, true)
	)

	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "a")
	logger.Error(NewPrefix("x", "y", "z"), "> %s", "b")
	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "c")
	processor.Shutdown()

	outLines := []string{
		"2018-08-26 09:29:03.123 | x/y/z: > a\n",
		"2018-08-26 09:29:03.123 | x/y/z: > c\n",
	}

	errLines := []string{
		"2018-08-26 09:29:03.123 | x/y/z: > b\n",
	}

	Expect(outfile.String()).To(Equal(strings.Join(outLines, "")))
	Expect(errfile.String()).To(Equal(strings.Join(errLines, "")))
}

func (s *ProcessorSuite) TestNonVerbose(t sweet.T) {
	clock := glock.NewMockClock()
	processor := newProcessor(false, false, clock)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outfile = &logWriter{}
		errfile = &logWriter{}
		logger  = processor.Logger(outfile, errfile, true)
	)

	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "a")
	logger.Info(NewPrefix("x", "y", "z"), "> %s", "b")
	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "c")
	processor.Shutdown()

	lines := []string{
		"2018-08-26 09:29:03.123 | x/y/z: > b\n",
	}

	Expect(outfile.String()).To(Equal(strings.Join(lines, "")))
	Expect(errfile.String()).To(BeEmpty())
}

//
// Helpers

type logWriter struct {
	bytes.Buffer
}

func (w *logWriter) String() string {
	return string(w.Bytes())
}

func (w *logWriter) Close() error {
	return nil
}
