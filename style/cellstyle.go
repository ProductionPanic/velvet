package style

type CellStyle struct {
	Fg Color // Foreground color
	Bg Color // Background color

	Bold      bool
	Italic    bool
	Underline bool
}

func DefaultCellStyle() CellStyle {
	return CellStyle{}
}
