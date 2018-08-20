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
		processor *processor
		outfile   io.Writer
		errfile   io.Writer
	}
)

func newLogger(processor *processor, outfile, errfile io.Writer) Logger {
	return &logger{
		processor: processor,
		outfile:   outfile,
		errfile:   errfile,
	}
}

func (l *logger) Debug(prefix *Prefix, format string, args ...interface{}) {
	if !l.processor.verbose {
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
	stream, file := l.getTargets(level)

	l.processor.enqueue(&message{
		level:     level,
		format:    format,
		args:      args,
		timestamp: time.Now(),
		prefix:    prefix,
		stream:    stream,
		file:      file,
	})
}

func (l *logger) getTargets(level LogLevel) (io.Writer, io.Writer) {
	if level == LevelError {
		return os.Stderr, l.errfile
	}

	return os.Stdout, l.outfile
}
