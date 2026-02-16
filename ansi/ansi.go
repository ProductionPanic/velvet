package ansi

import (
	"fmt"
	"strconv"

	"github.com/ProductionPanic/velvet/style/color"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	EnterAltScreen = "\033[?1049h"
	ExitAltScreen  = "\033[?1049l"

	ClearScreen = "\033[2J"

	ShowCursor       = "\033[?25h"
	HideCursor       = "\033[?25l"
	ClearLine        = "\033[2K" // Clear entire current line
	ClearLineToRight = "\033[K"  // Clear from cursor to end of line
	ClearLineToLeft  = "\033[1K" // Clear from beginning of line to cursor
	ClearBelow       = "\033[J"  // Clear everything below the cursor
	SaveCursor       = "\033[s"
	RestoreCursor    = "\033[u"
)

func SetCursorPosition(x, y int) string {
	return "\033[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H"
}

func ForegroundColor(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// 38;2; indicates TrueColor (RGB) foreground
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func BackgroundColor(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// 48;2; indicates TrueColor (RGB) background
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

func SetTitle(title string) string {
	return "\033]0;" + title + "\007"
}

func Print(s ...string) {
	fmt.Print(Combine(s...))
}

func Combine(s ...string) string {
	var result string
	for _, str := range s {
		result += str
	}
	return result
}

func ResetAll() {
	Print(
		ExitAltScreen,
		ClearScreen,
		SetCursorPosition(0, 0),
		ShowCursor,
		Reset,
	)
}

func MoveUp(n int) string    { return fmt.Sprintf("\033[%dA", n) }
func MoveDown(n int) string  { return fmt.Sprintf("\033[%dB", n) }
func MoveRight(n int) string { return fmt.Sprintf("\033[%dC", n) }
func MoveLeft(n int) string  { return fmt.Sprintf("\033[%dD", n) }
