package buffer

type Grid struct {
	Cells  []Cell
	Width  int
	Height int
}

func (g *Grid) index(x, y int) int {
	return y*g.Width + x
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

func (g *Grid) Clear() {
	for i := range g.Cells {
		g.Cells[i] = EmptyCell()
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
