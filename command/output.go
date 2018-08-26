package command

import (
	"bytes"

	"github.com/efritz/ij/logging"
)

type (
	outputProcessor interface {
		Process(line string)
	}

	outputProcessorFunc func(line string)

	stringProcessor struct {
		buffer bytes.Buffer
	}
)

var nilProcessor = outputProcessorFunc(func(_ string) {})

func (f outputProcessorFunc) Process(line string) {
	f(line)
}

func newLogProcessor(prefix *logging.Prefix, logFunc logging.LogFunc) outputProcessor {
	return outputProcessorFunc(func(line string) { logFunc(prefix, line) })
}

func newStringProcessor() *stringProcessor {
	return &stringProcessor{}
}

func (p *stringProcessor) Process(line string) {
	p.buffer.WriteString(line)
	p.buffer.WriteString("\n")
}

func (p *stringProcessor) String() string {
	return p.buffer.String()
}
