package logging

import (
	"fmt"
	"io"
	"time"
)

type (
	message struct {
		level       LogLevel
		format      string
		args        []interface{}
		timestamp   time.Time
		prefix      *Prefix
		writePrefix bool
		stream      io.Writer
		file        io.Writer
	}

	LogLevel int
)

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
)

func (m *message) Text() string {
	return fmt.Sprintf(m.format, m.args...)
}
