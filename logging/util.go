package logging

import (
	"fmt"
	"io"
	"os"
)

func emergencyLog(format string, args ...interface{}) {
	writeAll(os.Stderr, []byte(fmt.Sprintf(format, args...)))
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
