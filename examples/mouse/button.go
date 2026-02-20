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
	for i := range in.Cells {
		in.Cells[i].SetZoneId(b.ZoneID)
	}

	return in
}
