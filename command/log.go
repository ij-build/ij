package command

import "github.com/efritz/pvc/logging"

func newLogProcessor(logFunc logging.LogFunc) outputProcessor {
	return outputProcessorFunc(func(line string) { logFunc(line) })
}
