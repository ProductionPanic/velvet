package main

import (
	"fmt"
	"log"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
)

type model struct {
	button                 ButtonComponent
	offsetX, offsetY       int
	dragStartX, dragStartY int
	isDragging             bool
	width, height          int
	clickLog               string
}

func initialModel() velvet.Model {
	btn := ButtonComponent{
		ZoneID: "button1",
		Label:  "Drag me!",
	}

	return model{
		button:   btn,
		offsetX:  10,
		offsetY:  5,
		clickLog: "Click the button and drag it around!",
	}
}

func main() {
	f, err := velvet.LogToFile("mouse_events.log", "debug")
	if err != nil {
		panic(fmt.Sprintf("Failed to set up logging: %v", err))
	}
	defer f.Close()

	p := velvet.NewProgram(initialModel(), velvet.WithMouseAllMotion(true))
	if err := p.Run(); err != nil {
		panic(err)
	}
}

func (m model) Init() velvet.Cmd {
	return nil
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case velvet.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	case velvet.MouseMsg:
		action := "clicked"
		if msg.Type == velvet.MouseMotion {
			action = "moved"
		} else if msg.Type == velvet.MouseRelease {
			action = "released"
		}
		m.clickLog = fmt.Sprintf("%s at x:%d and y:%d, zone: %s", action, msg.X, msg.Y, msg.ZoneID)
		m.button, _ = m.button.Update(msg)
		return m, nil
	}
	return m, nil
}

func debug(grid *buffer.Grid) {
	zoneIds := []string{}

	for _, cell := range grid.Cells {
		if cell.ZoneId != "" {
			zoneIds = append(zoneIds, cell.ZoneId)
		}
	}

	log.Println("Zones in current render:", zoneIds)
}

func (m model) Render() velvet.Drawable {
	d := velvet.NewFlexContainer().SetSize(m.width, m.height)

	clicklog := velvet.NewComponent().Content(m.clickLog).Border(style.RoundedBorder())
	d.Add(m.button.Render(), clicklog)

	return d
}
