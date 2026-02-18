package velvet

import "github.com/ProductionPanic/velvet/buffer"

type FlexContainer struct {
	items         []*Container // items to be laid out
	width, height int          // the available space for the container
	horizontal    bool         // whether to layout items horizontally (default) or vertically
}

func NewFlexContainer() *FlexContainer {
	return &FlexContainer{
		items:      []*Container{},
		horizontal: true,
	}
}

func (fc *FlexContainer) SetWidth(w int) *FlexContainer {
	fc.width = w
	return fc
}

func (fc *FlexContainer) SetHeight(h int) *FlexContainer {
	fc.height = h
	return fc
}

func (fc *FlexContainer) SetSize(w, h int) *FlexContainer {
	fc.width = w
	fc.height = h
	return fc
}

func (fc *FlexContainer) Add(item ...*Container) *FlexContainer {
	fc.items = append(fc.items, item...)
	return fc
}

func (fc *FlexContainer) Horizontal() *FlexContainer {
	fc.horizontal = true
	return fc
}

func (fc *FlexContainer) Vertical() *FlexContainer {
	fc.horizontal = false
	return fc
}

func (fc *FlexContainer) Render() *buffer.Grid {
	out := buffer.NewGrid(fc.width, fc.height)

	if len(fc.items) == 0 {
		return out
	}

	if fc.horizontal {
		itemWidth := fc.width / len(fc.items)
		leftover := fc.width % len(fc.items)
		gap := 0
		if len(fc.items) > 1 {
			gap = leftover / (len(fc.items) - 1)
		}

		x := 0
		for i, item := range fc.items {
			item.Width(itemWidth).Height(fc.height)
			itemGrid := item.Render()
			out.Place(itemGrid, buffer.PositionX(x), buffer.AlignTop)
			x += itemWidth
			if i < len(fc.items)-1 {
				x += gap
			}
		}
	} else {
		itemHeight := fc.height / len(fc.items)
		leftover := fc.height % len(fc.items)
		gap := 0
		if len(fc.items) > 1 {
			gap = leftover / (len(fc.items) - 1)
		}

		y := 0
		for i, item := range fc.items {
			item.Width(fc.width).Height(itemHeight)
			itemGrid := item.Render()
			out.Place(itemGrid, buffer.AlignLeft, buffer.PositionY(y))
			y += itemHeight
			if i < len(fc.items)-1 {
				y += gap
			}
		}
	}

	return out
}

type GridContainer struct {
	items         []*Container // items to be laid out
	width, height int          // the available space for the container
	Rows, Cols    int          // number of rows and columns in the grid
}

func NewGridContainer() *GridContainer {
	return &GridContainer{
		items: []*Container{},
	}
}

func (gc *GridContainer) SetWidth(w int) *GridContainer {
	gc.width = w
	return gc
}

func (gc *GridContainer) SetHeight(h int) *GridContainer {
	gc.height = h
	return gc
}

func (gc *GridContainer) SetSize(w, h int) *GridContainer {
	gc.width = w
	gc.height = h
	return gc
}

func (gc *GridContainer) SetRows(rows int) *GridContainer {
	gc.Rows = rows
	return gc
}

func (gc *GridContainer) SetCols(cols int) *GridContainer {
	gc.Cols = cols
	return gc
}

func (gc *GridContainer) Add(item ...*Container) *GridContainer {
	gc.items = append(gc.items, item...)
	return gc
}

func (gc *GridContainer) Render() *buffer.Grid {
	out := buffer.NewGrid(gc.width, gc.height)

	if len(gc.items) == 0 {
		return out
	}

	if gc.Rows == 0 && gc.Cols == 0 {
		return out
	}

	// if only rows or cols is set, calculate the other based on the number of items
	if gc.Rows == 0 {
		gc.Rows = (len(gc.items) + gc.Cols - 1) / gc.Cols
	} else if gc.Cols == 0 {
		gc.Cols = (len(gc.items) + gc.Rows - 1) / gc.Rows
	}

	cellWidth := gc.width / gc.Cols
	cellHeight := gc.height / gc.Rows

	leftover := gc.width % gc.Rows
	gapX := 0
	if gc.Cols > 1 {
		gapX = leftover / (gc.Cols - 1)
	}

	leftover = gc.height % gc.Rows
	gapY := 0
	if gc.Rows > 1 {
		gapY = leftover / (gc.Rows - 1)
	}

	for i, item := range gc.items {
		row := i / gc.Cols
		col := i % gc.Cols

		item.Width(cellWidth).Height(cellHeight)
		itemGrid := item.Render()
		out.Place(itemGrid, buffer.PositionX(col*(cellWidth+gapX)), buffer.PositionY(row*(cellHeight+gapY)))
	}

	return out
}
