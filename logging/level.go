package logging

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
)

func (l LogLevel) Prefix() string {
	switch l {
	case LevelDebug:
		return "[D]"
	case LevelInfo:
		return "[I]"
	case LevelError:
		return "[E]"
	}

	return ""
}
