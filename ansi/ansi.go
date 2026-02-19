package ansi

import (
	"fmt"
	"strconv"

	"github.com/ProductionPanic/velvet/style/color"
	"github.com/ProductionPanic/velvet/terminal"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	EnterAltScreen = "\033[?1049h"
	ExitAltScreen  = "\033[?1049l"

	EnableMouseAllMotion = "\033[?1003h"
	EnableMouseClick     = "\033[?1000h"
	EnableMouseSGRMouse  = "\033[?1006h"

	DisableMouseAllMotion = "\033[?1003l"
	DisableMouseClick     = "\033[?1000l"
	DisableMouseSGRMouse  = "\033[?1006l"

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

// ForegroundColor256 outputs a 256-color foreground color code
func ForegroundColor256(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// 38;5;n indicates 256-color foreground
	idx := color.ToNearest256(c)
	return fmt.Sprintf("\x1b[38;5;%dm", idx)
}

// BackgroundColor256 outputs a 256-color background color code
func BackgroundColor256(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// 48;5;n indicates 256-color background
	idx := color.ToNearest256(c)
	return fmt.Sprintf("\x1b[48;5;%dm", idx)
}

// ForegroundColor16 outputs a 16-color ANSI foreground color code
func ForegroundColor16(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// ANSI 16 colors: indices 0-7 are normal, 8-15 are bright
	idx := color.ToNearest16(c)
	if idx < 8 {
		return fmt.Sprintf("\x1b[3%dm", idx)
	}
	return fmt.Sprintf("\x1b[9%dm", idx-8)
}

// BackgroundColor16 outputs a 16-color ANSI background color code
func BackgroundColor16(c color.Color) string {
	if c == (color.Color{}) {
		return ""
	}
	// ANSI 16 colors: indices 0-7 are normal, 8-15 are bright
	idx := color.ToNearest16(c)
	if idx < 8 {
		return fmt.Sprintf("\x1b[4%dm", idx)
	}
	return fmt.Sprintf("\x1b[10%dm", idx-8)
}

// ForegroundColorAuto outputs a foreground color using the terminal's detected capability
func ForegroundColorAuto(c color.Color, capability terminal.ColorCapability) string {
	switch capability {
	case terminal.ColorTrueColor:
		return ForegroundColor(c)
	case terminal.Color256:
		return ForegroundColor256(c)
	case terminal.Color16:
		return ForegroundColor16(c)
	default:
		return ForegroundColor(c)
	}
}

// BackgroundColorAuto outputs a background color using the terminal's detected capability
func BackgroundColorAuto(c color.Color, capability terminal.ColorCapability) string {
	switch capability {
	case terminal.ColorTrueColor:
		return BackgroundColor(c)
	case terminal.Color256:
		return BackgroundColor256(c)
	case terminal.Color16:
		return BackgroundColor16(c)
	default:
		return BackgroundColor(c)
	}
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
