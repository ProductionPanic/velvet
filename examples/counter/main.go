package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

type model struct {
	width, height int
	counter       int
	lastKey       string
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
		m.lastKey = fmt.Sprintf("%v", msg)

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, func() velvet.Msg {
				return velvet.Quit()
			}
		case "up":
			m.counter++
		case "down":
			m.counter--
		case " ":
			// Reset counter on space
			m.counter = 0
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	// Create a container
	container := velvet.NewComponent().
		Foreground(color.FromHex("#ffffff")).
		Background(color.FromHex("#111")).
		Padding(1, 2).
		Border(style.DoubleBorder()).
		BorderColor(color.FromHex("#00ff00")).
		Width(m.width / 2)

	displayText := fmt.Sprintf(
		"[bold,yellow]Counter Demo[]\n\n"+
			"counter value: [cyan,bold]%d[]\n"+
			"[dim]Press SPACE to reset\n"+
			"Press ESC or Ctrl+C to quit[]",
		m.counter,
	)

	rendered := container.Render(displayText)
	b.Place(rendered, buffer.AlignCenter, buffer.AlignMiddle)

	return b
}

func main() {
	p := velvet.NewProgram(model{})
	if err := p.Run(); err != nil {
		panic(err)
	}
}
