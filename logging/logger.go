package logging

import (
	"io"
	"os"
	"time"
)

type (
	Logger interface {
		Debug(prefix *Prefix, format string, args ...interface{})
		Info(prefix *Prefix, format string, args ...interface{})
		Error(prefix *Prefix, format string, args ...interface{})
	}

	LogFunc func(prefix *Prefix, format string, args ...interface{})

	logger struct {
		processor   *processor
		outfile     io.Writer
		errfile     io.Writer
		writePrefix bool
	}

	nilLogger struct{}
)

var NilLogger = &nilLogger{}

func newLogger(processor *processor, outfile, errfile io.Writer, writePrefix bool) Logger {
	return &logger{
		processor:   processor,
		outfile:     outfile,
		errfile:     errfile,
		writePrefix: writePrefix,
	}
}

func (l *logger) Debug(prefix *Prefix, format string, args ...interface{}) {
	if l == nil || !l.processor.verbose {
		return
	}

	l.enqueue(LevelDebug, prefix, format, args)
}

func (l *logger) Info(prefix *Prefix, format string, args ...interface{}) {
	l.enqueue(LevelInfo, prefix, format, args)
}

func (l *logger) Error(prefix *Prefix, format string, args ...interface{}) {
	l.enqueue(LevelError, prefix, format, args)
}

func (l *logger) enqueue(level LogLevel, prefix *Prefix, format string, args []interface{}) {
	if l == nil {
		return
	}

	stream, file := l.getTargets(level)

	l.processor.enqueue(&message{
		level:       level,
		format:      format,
		args:        args,
		timestamp:   time.Now(),
		prefix:      prefix,
		writePrefix: l.writePrefix,
		stream:      stream,
		file:        file,
	})
}

func (l *logger) getTargets(level LogLevel) (io.Writer, io.Writer) {
	if level == LevelError {
		return os.Stderr, l.errfile
	}

	return os.Stdout, l.outfile
}

func (l *nilLogger) Debug(*Prefix, string, ...interface{}) {}
func (l *nilLogger) Info(*Prefix, string, ...interface{})  {}
func (l *nilLogger) Error(*Prefix, string, ...interface{}) {}
