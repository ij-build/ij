package logging

import "github.com/mgutz/ansi"

type colorPicker struct {
	index int
}

var colors = []string{
	ansi.Black,
	ansi.Red,
	ansi.Green,
	ansi.Yellow,
	ansi.Blue,
	ansi.Magenta,
	ansi.Cyan,
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

func newColorPicker() *colorPicker {
	return &colorPicker{}
}

func (cp *colorPicker) next() string {
	index := cp.index % len(colors)
	cp.index++
	return colors[index]
}
