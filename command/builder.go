package command

type (
	Builder struct {
		builders []BuildFunc
		args     []string
	}

	BuildFunc func(b *Builder) error
)

func NewBuilder(builders []BuildFunc, args []string) *Builder {
	return &Builder{
		builders: builders,
		args:     args,
	}
}

func (b *Builder) Build() ([]string, error) {
	for _, builder := range b.builders {
		if err := builder(b); err != nil {
			return nil, err
		}
	}

	return b.args, nil
}

func (b *Builder) AddFlag(flag string) {
	b.args = append(b.args, flag)
}

func (b *Builder) AddFlagValue(flag, value string) {
	if value != "" {
		b.args = append(b.args, flag, value)
	}
}
