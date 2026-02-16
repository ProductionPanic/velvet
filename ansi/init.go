package ansi

import (
	"strconv"

	"github.com/ProductionPanic/velvet/style"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	EnterAltScreen = "\033[?1049h"
	ExitAltScreen  = "\033[?1049l"

	ShowCursor = "\033[?25h"
	HideCursor = "\033[?25l"
)

func SetCursorPosition(x, y int) string {
	return "\033[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H"
}

func SetForegroundColor(color style.Color) string {
	// use r g b to set foreground color ansi
	return "\033[38;2;" + strconv.Itoa(int(color.R)) + ";" + strconv.Itoa(int(color.G)) + ";" + strconv.Itoa(int(color.B)) + "m"
}

func SetBackgroundColor(color style.Color) string {
	// use r g b to set background color ansi
	return "\033[48;2;" + strconv.Itoa(int(color.R)) + ";" + strconv.Itoa(int(color.G)) + ";" + strconv.Itoa(int(color.B)) + "m"
}

func RenderCellStyle(cellStyle style.CellStyle) string {
	var s string
	if cellStyle.Bold {
		s += Bold
	}
	if cellStyle.Italic {
		s += Italic
	}
	if cellStyle.Underline {
		s += Underline
	}
	if cellStyle.Fg != (style.Color{}) {
		s += SetForegroundColor(cellStyle.Fg)
	}
	if cellStyle.Bg != (style.Color{}) {
		s += SetBackgroundColor(cellStyle.Bg)
	}
	return s
}

func SetTitle(title string) string {
	return "\033]0;" + title + "\007"
}
