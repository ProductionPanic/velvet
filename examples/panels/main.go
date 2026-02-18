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
	leftCounter   int
	rightCounter  int
	isLeftActive  bool
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
		switch msg.Type {
		case velvet.KeyTab:
			m.isLeftActive = !m.isLeftActive
		case velvet.KeyCtrlC, velvet.KeyEsc:
			return m, velvet.Quit
		case velvet.KeyUp:
			if m.isLeftActive {
				m.leftCounter++
			} else {
				m.rightCounter++
			}
		case velvet.KeyDown:
			if m.isLeftActive {
				m.leftCounter--
			} else {
				m.rightCounter--
			}
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	defaultPanel := velvet.NewContainer().
		Width(m.width/2).
		Height(m.height).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#555")).
		Padding(0, 1).
		Background(color.FromHex("#222")).
		Foreground(color.FromHex("#fff"))

	left := defaultPanel.Copy().BorderTextTop("Counter one")
	right := defaultPanel.Copy().BorderTextTop("Counter two")

	if m.isLeftActive {
		left.BorderColor(color.FromHex("#ff00ff"))
	} else {
		right.BorderColor(color.FromHex("#ff00ff"))
	}

	leftPanelText := fmt.Sprintf(
		"[bold,magenta]Left Panel[]\n\n"+
			"counter value: [cyan,bold]%d[]\n"+
			"Press ESC or Ctrl+C to quit[]",
		m.leftCounter,
	)

	rightPanelText := fmt.Sprintf(
		"[bold,magenta]Right Panel[]\n\n"+
			"counter value: [cyan,bold]%d[]\n",
		m.rightCounter,
	)

	leftPanel := left.Render(leftPanelText)
	rightPanel := right.Render(rightPanelText)

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
