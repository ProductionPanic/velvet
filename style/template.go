package style

import (
	"strings"

	"github.com/ProductionPanic/velvet/style/color"
)

func ApplyTemplateStyle(parts []string, style *CellStyle) *CellStyle {

	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "bold":
			style.Bold = true
		case "italic":
			style.Italic = true
		case "underline":
			style.Underline = true
		default:
			if part == "" {
				continue
			}
			// check if it starts with "bg" for background color
			if len(part) > 2 && part[:2] == "bg" {
				style.Bg = colorFromString(part[2:])
			} else {
				style.Fg = colorFromString(part)
			}
		}
	}

	return style
}

func colorFromString(s string) color.Color {
	// check if it's a hex color
	if len(s) == 7 && s[0] == '#' {
		return color.FromHex(s)
	}

	// check for named colors
	switch s {
	case "black":
		return color.Black()
	case "red":
		return color.Red()
	case "green":
		return color.Green()
	case "yellow":
		return color.Yellow()
	case "blue":
		return color.Blue()
	case "magenta":
		return color.Magenta()
	case "cyan":
		return color.Cyan()
	case "white":
		return color.White()
	default:
		return color.Default() // default to magenta if unknown color
	}
}
