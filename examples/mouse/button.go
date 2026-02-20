package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style/color"
)

type ButtonComponent struct {
	Width   int
	Height  int
	ZoneID  string
	Label   string
	bgColor color.Color
}

func (b ButtonComponent) Update(msg velvet.Msg) (ButtonComponent, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.MouseMsg:
		if msg.Type == velvet.MouseMotion {
			if msg.ZoneID == b.ZoneID {
				b.bgColor = color.FromHex("#ddd") // Change background color on hover
			} else {
				b.bgColor = color.FromHex("#aaa") // Default background color
			}
		}
	}
	return b, nil
}

func (b ButtonComponent) Render() velvet.Drawable {
	component := velvet.NewComponent().
		Width(b.Width).
		Height(b.Height).
		Background(b.bgColor).
		Foreground(color.Black()).
		Content(b.Label)

	if b.ZoneID != "" {
		component.SetZoneID(b.ZoneID)
	}

	return component
}
