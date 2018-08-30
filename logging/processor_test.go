package logging

import (
	"bytes"
	"strings"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type ProcessorSuite struct{}

func (s *ProcessorSuite) TestBasic(t sweet.T) {
	var (
		clock     = glock.NewMockClock()
		outStream = &bytes.Buffer{}
		errStream = &bytes.Buffer{}
	)

	processor := newProcessor(true, false, clock, outStream, errStream)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outFile = &logWriter{}
		errFile = &logWriter{}
		logger  = processor.Logger(outFile, errFile, true)
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

	Expect(outFile.String()).To(Equal(strings.Join(lines, "")))
	Expect(errFile.String()).To(BeEmpty())
}

func (s *ProcessorSuite) TestErrors(t sweet.T) {
	var (
		clock     = glock.NewMockClock()
		outStream = &bytes.Buffer{}
		errStream = &bytes.Buffer{}
	)

	processor := newProcessor(true, false, clock, outStream, errStream)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outFile = &logWriter{}
		errFile = &logWriter{}
		logger  = processor.Logger(outFile, errFile, true)
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

	Expect(outFile.String()).To(Equal(strings.Join(outLines, "")))
	Expect(errFile.String()).To(Equal(strings.Join(errLines, "")))
}

func (s *ProcessorSuite) TestNonVerbose(t sweet.T) {

	var (
		clock     = glock.NewMockClock()
		outStream = &bytes.Buffer{}
		errStream = &bytes.Buffer{}
	)

	processor := newProcessor(false, false, clock, outStream, errStream)
	processor.Start()

	clock.SetCurrent(time.Unix(1535293743, 123000000))

	var (
		outFile = &logWriter{}
		errFile = &logWriter{}
		logger  = processor.Logger(outFile, errFile, true)
	)

	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "a")
	logger.Info(NewPrefix("x", "y", "z"), "> %s", "b")
	logger.Debug(NewPrefix("x", "y", "z"), "> %s", "c")
	processor.Shutdown()

	lines := []string{
		"2018-08-26 09:29:03.123 | x/y/z: > b\n",
	}

	Expect(outFile.String()).To(Equal(strings.Join(lines, "")))
	Expect(errFile.String()).To(BeEmpty())
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
