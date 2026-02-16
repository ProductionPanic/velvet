package buffer

import "github.com/ProductionPanic/velvet/style"

type Cell struct {
	Rune  rune
	Style *style.CellStyle
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
