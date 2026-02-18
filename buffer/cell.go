package buffer

import (
	"github.com/ProductionPanic/velvet/style"
	"github.com/mattn/go-runewidth"
)

type Cell struct {
	Content   string // changed from rune to string to support multi-rune grapheme clusters
	Style     *style.CellStyle
	runeWidth int
}

func NewCell(r rune, s *style.CellStyle) Cell {
	return Cell{
		Content: string(r),
		Style:   s,
	}
}

func EmptyCell() Cell {
	return Cell{
		Content: " ",
		Style:   style.DefaultCellStyle(),
	}
}

func (c Cell) Width() int {
	if c.runeWidth == 0 {
		c.runeWidth = runewidth.StringWidth(c.Content)
	}
	return c.runeWidth
}

func (c Cell) IsSpace() bool {
	return c.Content == " "
}
