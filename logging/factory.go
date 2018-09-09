package logging

import "io"

type (
	LoggerFactory struct {
		logProcessor Processor
		fileFactory  FileFactory
	}

	FileFactory func(prefix string) (io.WriteCloser, io.WriteCloser, error)
)

func NewLoggerFactory(logProcessor Processor, fileFactory FileFactory) *LoggerFactory {
	return &LoggerFactory{
		logProcessor: logProcessor,
		fileFactory:  fileFactory,
	}
}

func (f *LoggerFactory) Logger(prefix string, writePrefix bool) (Logger, error) {
	outFile, errFile, err := f.fileFactory(prefix)
	if err != nil {
		return nil, err
	}

	return f.logProcessor.Logger(outFile, errFile, writePrefix), nil
}
