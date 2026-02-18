package buffer

import (
	"github.com/ProductionPanic/velvet/style"
	"github.com/mattn/go-runewidth"
)

type Cell struct {
	Rune      rune
	Style     *style.CellStyle
	runeWidth int
}

func NewCell(r rune, s *style.CellStyle) Cell {
	return Cell{
		Rune:  r,
		Style: s,
	}
}

func EmptyCell() Cell {
	return Cell{
		Rune:  ' ',
		Style: style.DefaultCellStyle(),
	}
}

func (c Cell) Width() int {
	if c.runeWidth == 0 {
		c.runeWidth = runewidth.RuneWidth(c.Rune)
	}
	return c.runeWidth
}

func (c Cell) IsSpace() bool {
	return c.Rune == ' '
}
