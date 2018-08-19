package logging

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
)
