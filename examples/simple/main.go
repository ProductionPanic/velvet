package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
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
		// Quit on Ctrl+C or Escape
		if msg.Type == velvet.KeyCtrlC || msg.Type == velvet.KeyEsc {
			return m, func() velvet.Msg {
				return velvet.Quit()
			}
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	component := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Padding(1).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#3300ff")).
		Width(m.width / 4)

	content := "Hello [red]velvet![] This is a demo of the Velvet [bold,cyan]TUI library[]. " +
		"It supports text wrapping, borders, padding, and more!\n\n" +
		"[dim]Press ESC or Ctrl+C to quit[]"

	b.Place(component.Render(content), buffer.AlignCenter, buffer.AlignMiddle)

	return b
}

func main() {
	// That's it! Much simpler than the original demo
	p := velvet.NewProgram(model{})
	if err := p.Run(); err != nil {
		panic(err)
	}
}
