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

		switch msg.Type {
		case velvet.KeyCtrlC, velvet.KeyEsc:
			return m, func() velvet.Msg {
				return velvet.Quit()
			}
		case velvet.KeyUp:
			m.counter++
		case velvet.KeyDown:
			m.counter--
		case velvet.KeySpace:
			// Reset counter on space
			m.counter = 0
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	panel := velvet.NewContainer().
		Width(m.width/2).
		Height(m.height).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#ff00ff")).
		Padding(0, 1).
		Background(color.FromHex("#222")).
		Foreground(color.FromHex("#fff"))

	leftPanelText := fmt.Sprintf(
		"[bold,magenta]Left Panel[]\n\n"+
			"counter value: [cyan,bold]%d[]\n"+
			"[dim]Press SPACE to reset\n"+
			"Press ESC or Ctrl+C to quit[]",
		m.counter,
	)

	rightPanelText := fmt.Sprintf(
		"[bold,magenta]Right Panel[]\n\n"+
			"Last key: [yellow,bold]%s[]\n",
		m.lastKey,
	)

	leftPanel := panel.BorderTextTop("Counter").Render(leftPanelText)
	rightPanel := panel.BorderTextTop("Keylogger").Render(rightPanelText)

	b.Place(leftPanel, buffer.AlignLeft, buffer.AlignTop)
	b.Place(rightPanel, buffer.AlignRight, buffer.AlignTop)

	return b
}

func main() {
	p := velvet.NewProgram(model{})
	if err := p.Run(); err != nil {
		panic(err)
	}
}
