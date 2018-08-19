package logging

import (
	"fmt"
	"io"
	"time"
)

type message struct {
	level     LogLevel
	format    string
	args      []interface{}
	timestamp time.Time
	prefix    string
	colorCode string
	stream    io.Writer
	file      io.Writer
}

func (m *message) Text() string {
	return fmt.Sprintf(m.format, m.args...)
}
