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

	if fc.horizontal {
		itemWidth := fc.width / len(fc.items)
		for i, item := range fc.items {
			item.Width(itemWidth).Height(fc.height)
			itemGrid := item.Render()
			out.Place(itemGrid, buffer.PositionX(i*itemWidth), buffer.AlignTop)
		}
	} else {
		itemHeight := fc.height / len(fc.items)
		for i, item := range fc.items {
			item.Width(fc.width).Height(itemHeight)
			itemGrid := item.Render()
			out.Place(itemGrid, buffer.AlignLeft, buffer.PositionY(i*itemHeight))
		}
	}

	return out
}
