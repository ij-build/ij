package command

type (
	outputProcessor interface {
		Process(line string)
	}

	outputProcessorFunc func(line string)
)

var nilProcessor = outputProcessorFunc(func(_ string) {})

func (f outputProcessorFunc) Process(line string) {
	f(line)
}
