package logging

import (
	"fmt"
	"io"
	"os"
)

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
