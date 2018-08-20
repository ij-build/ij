package logging

import (
	"io"
	"os"
	"time"
)

type (
	Logger interface {
		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	LogFunc func(format string, args ...interface{})

	logger struct {
		processor *processor
		prefix    string
		colorCode string
		outfile   io.Writer
		errfile   io.Writer
	}
)

func newLogger(processor *processor, prefix, colorCode string, outfile, errfile io.Writer) Logger {
	return &logger{
		processor: processor,
		prefix:    prefix,
		colorCode: colorCode,
		outfile:   outfile,
		errfile:   errfile,
	}
}

func (l *logger) Debug(format string, args ...interface{}) {
	if !l.processor.verbose {
		return
	}

	l.enqueue(LevelDebug, format, args)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.enqueue(LevelInfo, format, args)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.enqueue(LevelError, format, args)
}

func (l *logger) enqueue(level LogLevel, format string, args []interface{}) {
	l.processor.enqueue(&message{
		level:     level,
		format:    format,
		args:      args,
		timestamp: time.Now(),
		prefix:    l.prefix,
		colorCode: l.colorCode,
		stream:    getStream(level),
		file:      l.outfile,
	})
}

func getStream(level LogLevel) io.Writer {
	if level == LevelError {
		return os.Stderr
	}

	return os.Stdout
}
