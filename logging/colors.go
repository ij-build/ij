package logging

import "github.com/mgutz/ansi"

var colors = []string{
	ansi.Red,
	ansi.Yellow,
	ansi.Blue,
	ansi.Magenta,
	ansi.White,
	ansi.LightBlack,
	ansi.LightRed,
	ansi.LightGreen,
	ansi.LightYellow,
	ansi.LightBlue,
	ansi.LightMagenta,
	ansi.LightCyan,
	ansi.LightWhite,
}

var levelColors = map[LogLevel]string{
	LevelDebug: ansi.Cyan,
	LevelInfo:  ansi.Green,
	LevelWarn:  ansi.Yellow,
	LevelError: ansi.ColorCode("red+b"),
}
