package velvet

import "github.com/ProductionPanic/velvet/style"

type Buffer struct {
	Width  int
	Height int
	Cells  []Cell // A flat array of cells, row-major order
}

func (b *Buffer) index(x, y int) int {
	return y*b.Width + x
}

func (b *Buffer) Clear() {
	for i := range b.Cells {
		b.Cells[i] = Cell{Char: ' ', Style: style.DefaultCellStyle()}
	}
}

func (b *Buffer) Crop(newWidth, newHeight int) *Buffer {
	cells := make([]Cell, newWidth*newHeight)
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			if x < b.Width && y < b.Height {
				cells[y*newWidth+x] = b.Cells[b.index(x, y)]
			} else {
				cells[y*newWidth+x] = Cell{Char: ' ', Style: style.DefaultCellStyle()}
			}
		}
	}
	return &Buffer{
		Width:  newWidth,
		Height: newHeight,
		Cells:  cells,
	}
}

func NewBuffer(width, height int) *Buffer {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = Cell{Char: ' ', Style: style.DefaultCellStyle()}
	}
	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}
