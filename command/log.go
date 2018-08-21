package command

import "github.com/efritz/ij/logging"

func newLogProcessor(prefix *logging.Prefix, logFunc logging.LogFunc) outputProcessor {
	return outputProcessorFunc(func(line string) { logFunc(prefix, line) })
}
