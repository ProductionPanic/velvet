package style

import "github.com/ProductionPanic/velvet/style/color"

type CellStyle struct {
	Fg color.Color // Foreground color
	Bg color.Color // Background color

	Bold      bool
	Italic    bool
	Underline bool
}

func DefaultCellStyle() *CellStyle {
	return &CellStyle{}
}

func (s *CellStyle) SetFg(r, g, b uint8) *CellStyle {
	s.Fg = color.FromRGB(r, g, b)
	return s
}

func (s *CellStyle) SetBg(r, g, b uint8) *CellStyle {
	s.Bg = color.FromRGB(r, g, b)
	return s
}

func (s *CellStyle) SetBold(b bool) *CellStyle {
	s.Bold = b
	return s
}

func (s *CellStyle) SetItalic(i bool) *CellStyle {
	s.Italic = i
	return s
}

func (s *CellStyle) SetUnderline(u bool) *CellStyle {
	s.Underline = u
	return s
}

func (s *CellStyle) Copy() *CellStyle {
	return &CellStyle{
		Fg:        s.Fg,
		Bg:        s.Bg,
		Bold:      s.Bold,
		Italic:    s.Italic,
		Underline: s.Underline,
	}
}

func NewCellStyle() *CellStyle {
	return DefaultCellStyle()

}
