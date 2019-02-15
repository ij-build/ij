package logging

import (
	"fmt"
	"os"
)

type nilWriter struct{}

var NilWriter = &nilWriter{}

func EmergencyLog(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}

func (w *nilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *nilWriter) Close() error {
	return nil
}
