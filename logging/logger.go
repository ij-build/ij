package logging

import "io"

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
		outfile   io.Writer
		errfile   io.Writer
	}
)

func newLogger(processor *processor, prefix string, outfile, errfile io.Writer) Logger {
	return &logger{
		processor: processor,
		prefix:    prefix,
		outfile:   outfile,
		errfile:   errfile,
	}
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.processor.enqueue(&message{LevelDebug, format, args, l.prefix, l.outfile})
}

func (l *logger) Info(format string, args ...interface{}) {
	l.processor.enqueue(&message{LevelInfo, format, args, l.prefix, l.outfile})
}

func (l *logger) Error(format string, args ...interface{}) {
	l.processor.enqueue(&message{LevelError, format, args, l.prefix, l.errfile})
}
