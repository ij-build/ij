package logging

import (
	"strings"
)

type Prefix struct {
	parts []string
}

func NewPrefix(parts ...string) *Prefix {
	return &Prefix{
		parts: parts,
	}
}

func (p *Prefix) Append(part string) *Prefix {
	parts := []string{}
	for _, part := range p.parts {
		parts = append(parts, part)
	}

	return NewPrefix(append(parts, part)...)
}

func (p *Prefix) Serialize(picker *colorPicker) string {
	colorized := []string{}
	for _, part := range p.parts {
		colorized = append(colorized, picker.colorize(part))
	}

	return strings.Join(colorized, "/")
}
