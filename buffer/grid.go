package buffer

import "github.com/ProductionPanic/velvet/style"

type Grid struct {
	Cells  []Cell
	Width  int
	Height int
}

func (g *Grid) index(x, y int) int {
	return y*g.Width + x
}

func (g *Grid) Coordinates(index int) (x, y int) {
	x = index % g.Width
	y = index / g.Width
	return x, y
}

func (g *Grid) isInBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

func (g *Grid) Set(x, y int, cell Cell) {
	if !g.isInBounds(x, y) {
		return
	}
	g.Cells[g.index(x, y)] = cell
}

func (g *Grid) Get(x, y int) Cell {
	if !g.isInBounds(x, y) {
		return EmptyCell()
	}
	return g.Cells[g.index(x, y)]
}

func (g *Grid) SetRune(x, y int, c string) {
	if !g.isInBounds(x, y) {
		return
	}
	cell := g.Get(x, y)
	cell.Content = c
	g.Set(x, y, cell)
}

func (g *Grid) SetEmptyCells(c string, s *style.CellStyle) {
	for i := range g.Cells {
		if g.Cells[i].Content == " " {
			g.Cells[i] = Cell{
				Content: c,
				Style:   s,
			}
		}
	}
}

// Resize just creates a new grid with the given dimensions, and copies over the cells that fit in the new dimensions. Cells that don't fit are discarded, and new cells are initialized as empty.
func (g *Grid) Resize(width, height int) *Grid {
	newGrid := NewGrid(width, height)
	for y := 0; y < height && y < g.Height; y++ {
		for x := 0; x < width && x < g.Width; x++ {
			newGrid.Set(x, y, g.Cells[g.index(x, y)])
		}
	}

	return newGrid
}

func (g *Grid) AddRow() {
	g.Cells = append(g.Cells, make([]Cell, g.Width)...)
	g.Height++
}

func (g *Grid) AddColumn() {
	g = g.Resize(g.Width+1, g.Height)
}

func (g *Grid) Clear() {
	for i := range g.Cells {
		g.Cells[i] = EmptyCell()
	}
}

func (g *Grid) Place(b *Grid, x PositionX, y PositionY) {
	var startX, startY int
	switch x {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (g.Width - b.Width) / 2
	case AlignRight:
		startX = g.Width - b.Width
	default:
		startX = int(x)
	}

	switch y {
	case AlignTop:
		startY = 0
	case AlignMiddle:
		startY = (g.Height - b.Height) / 2
	case AlignBottom:
		startY = g.Height - b.Height
	default:
		startY = int(y)
	}

	for by := 0; by < b.Height; by++ {
		for bx := 0; bx < b.Width; bx++ {
			g.Set(startX+bx, startY+by, b.Get(bx, by))
		}
	}
}

func NewGrid(width, height int) *Grid {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = EmptyCell()
	}

	return &Grid{
		Cells:  cells,
		Width:  width,
		Height: height,
	}
}
