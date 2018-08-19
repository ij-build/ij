package command

import "bytes"

type stringProcessor struct {
	buffer bytes.Buffer
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
