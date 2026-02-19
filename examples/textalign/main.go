package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

type model struct {
	width, height int
}

func (m model) Init() velvet.Cmd {
	return nil
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case velvet.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	}

	return m, nil
}

func (m model) Render() velvet.Drawable {
	// Left-aligned text
	leftBox := velvet.NewComponent().
		Content("[bold,blue]Left Aligned[]\n\nThis text is\naligned to the\nleft side.").
		TextAlign(velvet.AlignLeft).
		Width(m.width / 3).
		Height(8).
		Padding(1).
		Border(style.SimpleBorder()).
		BorderColor(color.FromHex("#0088ff")).
		Foreground(color.FromHex("#ffffff"))

	// Center-aligned text
	centerBox := velvet.NewComponent().
		Content("[bold,green]Center Aligned[]\n\nThis text is\ncentered in\nthe box.").
		TextAlign(velvet.AlignCenter).
		Width(m.width / 3).
		Height(8).
		Padding(1).
		Border(style.SimpleBorder()).
		BorderColor(color.FromHex("#00ff00")).
		Foreground(color.FromHex("#ffffff"))

	// Right-aligned text
	rightBox := velvet.NewComponent().
		Content("[bold,yellow]Right Aligned[]\n\nThis text is\naligned to the\nright side.").
		TextAlign(velvet.AlignRight).
		Width(m.width / 3).
		Height(8).
		Padding(1).
		Border(style.SimpleBorder()).
		BorderColor(color.FromHex("#ffff00")).
		Foreground(color.FromHex("#000000"))

	// Title
	title := velvet.NewComponent().
		Content("[bold,magenta]Text Alignment Demo[]\n\nUse TextAlign() to control horizontal text alignment in components.").
		Width(m.width).
		Height(3).
		Foreground(color.FromHex("#ff00ff"))

	// Create a flex layout to arrange the three boxes horizontally
	layout := velvet.NewFlexContainer().
		SetSize(m.width, m.height).
		Vertical().
		Add(
			title,
			velvet.NewFlexContainer().
				SetSize(m.width, 10).
				Horizontal().
				Add(leftBox, centerBox, rightBox),
		)

	return layout
}

func main() {
	p := velvet.NewProgram(model{})
	if err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
