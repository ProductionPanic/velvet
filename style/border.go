package style

type BorderStyle struct {
	Top    string
	Bottom string
	Left   string
	Right  string

	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
}

func RoundedBorder() BorderStyle {
	return BorderStyle{
		Top:      "─",
		Bottom:   "─",
		Left:     "│",
		Right:    "│",
		TopLeft:  "╭",
		TopRight: "╮",
	}
}
