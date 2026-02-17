package color

import (
	"fmt"
	"strconv"
	"strings"
)

type Color struct {
	R, G, B uint8
}

func FromRGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

func FromHex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		hex = fmt.Sprintf("%c%c%c%c%c%c", hex[0], hex[0], hex[1], hex[1], hex[2], hex[2])
	}

	if len(hex) != 6 {
		return Color{0, 0, 0} // default to black if invalid hex
	}

	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return Color{0, 0, 0} // default to black if parsing fails
	}

	return Color{
		R: uint8(value >> 16),
		G: uint8((value >> 8) & 0xFF),
		B: uint8(value & 0xFF),
	}
}
