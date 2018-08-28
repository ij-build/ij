package command

import "io"

type (
	Builder struct {
		prelude  []string
		options  []string
		args     []string
		stdin    io.ReadCloser
		builders []BuildFunc
	}

	BuildFunc func(b *Builder) error
)

func NewBuilder(prelude []string, builders []BuildFunc) *Builder {
	return &Builder{
		prelude:  prelude,
		builders: builders,
	}
}

func (b *Builder) Build() ([]string, io.ReadCloser, error) {
	for _, builder := range b.builders {
		if err := builder(b); err != nil {
			return nil, nil, err
		}
	}

	return append(b.prelude, append(b.options, b.args...)...), b.stdin, nil
}

func (b *Builder) AddArgs(args ...string) {
	b.args = append(b.args, args...)
}

func (b *Builder) AddFlag(flag string) {
	b.options = append(b.options, flag)
}

func (b *Builder) AddFlagValue(flag, value string) {
	if value != "" {
		b.options = append(b.options, flag, value)
	}
}

func (b *Builder) SetStdin(rc io.ReadCloser) {
	b.stdin = rc
}
