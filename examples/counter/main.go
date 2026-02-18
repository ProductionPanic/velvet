package main

import (
	"fmt"
	"time"

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
	// Return a command that ticks every second
	return velvet.Tick(time.Second)
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case velvet.TickMsg:
		m.counter++
		// Return another tick command to keep the counter going
		return m, velvet.Tick(time.Second)

	case velvet.KeyMsg:
		m.lastKey = fmt.Sprintf("%v", msg)

		switch msg.Type {
		case velvet.KeyCtrlC, velvet.KeyEsc:
			return m, func() velvet.Msg {
				return velvet.Quit()
			}
		case velvet.KeySpace:
			// Reset counter on space
			m.counter = 0
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	// Create a container
	container := velvet.NewContainer().
		Foreground(color.FromHex("#ffffff")).
		Padding(2).
		Border(style.DoubleBorder()).
		BorderColor(color.FromHex("#00ff00")).
		Width(m.width / 2)

	displayText := fmt.Sprintf(
		"[bold,yellow]Counter Demo[]\n\n"+
			"Seconds elapsed: [cyan,bold]%d[]\n"+
			"Last key pressed: [magenta]%s[]\n\n"+
			"[dim]Press SPACE to reset\n"+
			"Press ESC or Ctrl+C to quit[]",
		m.counter,
		m.lastKey,
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
