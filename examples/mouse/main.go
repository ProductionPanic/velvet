package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
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
	return model{
		button: ButtonComponent{
			ZoneID: "button1",
		},
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
		return m, nil
	}
	return m, nil
}

func (m model) Render() velvet.Drawable {
	d := velvet.NewFlexContainer().SetSize(m.width, m.height)

	clicklog := velvet.NewComponent().Content(m.clickLog).Border(style.RoundedBorder())

	d.Add(m.button, clicklog)

	return d
}
