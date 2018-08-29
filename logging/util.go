package logging

import (
	"fmt"
	"io"
	"os"
)

type nilWriter struct{}

var NilWriter = &nilWriter{}

func EmergencyLog(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if n > 0 {
			data = data[n:]
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *nilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *nilWriter) Close() error {
	return nil
}
