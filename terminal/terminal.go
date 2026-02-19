package terminal

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// ColorCapability represents the color support level of a terminal
type ColorCapability int

const (
	ColorTrueColor ColorCapability = iota // 24-bit RGB color (16.7 million colors)
	Color256                              // 256-color palette
	Color16                               // 16 ANSI colors
)

func GetSize() (width, height int, err error) {
	return term.GetSize(
		int(os.Stdout.Fd()),
	)
}

func RawMode() (restore func(), err error) {
	oldState, err := term.MakeRaw(
		int(os.Stdin.Fd()),
	)
	if err != nil {
		return nil, err
	}

	restore = func() {
		_ = term.Restore(
			int(os.Stdin.Fd()),
			oldState,
		)
	}

	return restore, nil
}

// DetectColorCapability detects the color capability of the terminal
// by checking environment variables
func DetectColorCapability() ColorCapability {
	colorterm := strings.ToLower(os.Getenv("COLORTERM"))
	if colorterm == "truecolor" || colorterm == "24bit" {
		return ColorTrueColor
	}

	termVar := strings.ToLower(os.Getenv("TERM"))

	// Check for 256-color support
	if strings.Contains(termVar, "256color") {
		return Color256
	}

	// Default to 16-color if TERM contains "color"
	if strings.Contains(termVar, "color") {
		return Color16
	}

	// Default fallback to 256-color for most terminals
	return Color256
}
