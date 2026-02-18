package color

func Black() Color {
	return Color{0, 0, 0}
}

func Red() Color {
	return Color{255, 0, 0}
}

func Green() Color {
	return Color{0, 255, 0}
}

func Yellow() Color {
	return Color{255, 255, 0}
}

func Blue() Color {
	return Color{0, 0, 255}
}

func Magenta() Color {
	return Color{255, 0, 255}
}

func Cyan() Color {
	return Color{0, 255, 255}
}

func White() Color {
	return Color{255, 255, 255}
}

func Default() Color {
	return Magenta()
}
