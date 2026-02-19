package color

import "math"

// Standard ANSI 16 color palette (basic + bright variants)
var ansi16Palette = [16]Color{
	{0, 0, 0},       // Black
	{128, 0, 0},     // Red
	{0, 128, 0},     // Green
	{128, 128, 0},   // Yellow
	{0, 0, 128},     // Blue
	{128, 0, 128},   // Magenta
	{0, 128, 128},   // Cyan
	{192, 192, 192}, // White
	{128, 128, 128}, // Bright Black
	{255, 0, 0},     // Bright Red
	{0, 255, 0},     // Bright Green
	{255, 255, 0},   // Bright Yellow
	{0, 0, 255},     // Bright Blue
	{255, 0, 255},   // Bright Magenta
	{0, 255, 255},   // Bright Cyan
	{255, 255, 255}, // Bright White
}

// xterm 256-color palette
// First 16 colors are ANSI colors, then 216-color cube (6x6x6), then 24 grayscale
var xterm256Palette = generateXterm256Palette()

func generateXterm256Palette() [256]Color {
	var palette [256]Color

	// First 16 colors are standard ANSI
	copy(palette[0:16], ansi16Palette[:])

	// 216-color cube (6x6x6)
	colorIndex := 16
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				palette[colorIndex] = Color{
					R: uint8((r * 255) / 5),
					G: uint8((g * 255) / 5),
					B: uint8((b * 255) / 5),
				}
				colorIndex++
			}
		}
	}

	// 24 grayscale colors
	for i := 0; i < 24; i++ {
		gray := uint8(8 + (i*247)/24)
		palette[colorIndex] = Color{R: gray, G: gray, B: gray}
		colorIndex++
	}

	return palette
}

// euclideanDistance calculates the Euclidean distance between two colors in RGB space
func euclideanDistance(c1, c2 Color) float64 {
	r := float64(int(c1.R) - int(c2.R))
	g := float64(int(c1.G) - int(c2.G))
	b := float64(int(c1.B) - int(c2.B))
	return math.Sqrt(r*r + g*g + b*b)
}

// ToNearest16 converts an RGB color to the nearest ANSI 16-color
func ToNearest16(c Color) uint8 {
	minDistance := math.MaxFloat64
	nearestIndex := uint8(0)

	for i, paletteColor := range ansi16Palette {
		distance := euclideanDistance(c, paletteColor)
		if distance < minDistance {
			minDistance = distance
			nearestIndex = uint8(i)
		}
	}

	return nearestIndex
}

// ToNearest256 converts an RGB color to the nearest xterm 256-color palette index
func ToNearest256(c Color) uint8 {
	minDistance := math.MaxFloat64
	nearestIndex := uint8(0)

	for i, paletteColor := range xterm256Palette {
		distance := euclideanDistance(c, paletteColor)
		if distance < minDistance {
			minDistance = distance
			nearestIndex = uint8(i)
		}
	}

	return nearestIndex
}
