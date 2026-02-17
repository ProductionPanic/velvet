package style

type BorderStyle struct {
	Top    rune
	Bottom rune
	Left   rune
	Right  rune

	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
}

func RoundedBorder() BorderStyle {
	return BorderStyle{
		Top:         '─',
		Bottom:      '─',
		Left:        '│',
		Right:       '│',
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
	}
}

func DoubleBorder() BorderStyle {
	return BorderStyle{
		Top:         '═',
		Bottom:      '═',
		Left:        '║',
		Right:       '║',
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
	}
}

func SimpleBorder() BorderStyle {
	return BorderStyle{
		Top:         '-',
		Bottom:      '-',
		Left:        '|',
		Right:       '|',
		TopLeft:     '+',
		TopRight:    '+',
		BottomLeft:  '+',
		BottomRight: '+',
	}
}
