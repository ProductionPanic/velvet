package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
)

type ButtonComponent struct {
	velvet.Component
	ZoneID string
	Label  string
}

func (b ButtonComponent) Render() *buffer.Grid {
	in := b.RenderWithContent()
	out := buffer.NewGrid(in.Width, in.Height)

	for i := 0; i < in.Width; i++ {
		for j := 0; j < in.Height; j++ {
			cell := in.Get(j, i)
			cell.SetZoneId(b.ZoneID)
			out.Set(j, i, cell)
		}
	}

	return out
}
