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
	content       string
	cursorPos     int
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
		case velvet.KeyCtrlC, velvet.KeyEsc:
			return m, func() velvet.Msg {
				return velvet.Quit()
			}

		case velvet.KeyEnter:
			m.content += "\n"
			m.cursorPos++

		case velvet.KeyBackspace:
			if len(m.content) > 0 && m.cursorPos > 0 {
				runes := []rune(m.content)
				m.content = string(runes[:m.cursorPos-1]) + string(runes[m.cursorPos:])
				m.cursorPos--
			}

		case velvet.KeyLeft:
			if m.cursorPos > 0 {
				m.cursorPos--
			}

		case velvet.KeyRight:
			if m.cursorPos < len([]rune(m.content)) {
				m.cursorPos++
			}

		case velvet.KeyRunes:
			runes := []rune(m.content)
			before := runes[:m.cursorPos]
			after := runes[m.cursorPos:]
			m.content = string(before) + string(msg.Runes) + string(after)
			m.cursorPos += len(msg.Runes)
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	// Create a container with border
	container := velvet.NewContainer().
		Foreground(color.FromHex("#ffffff")).
		Background(color.FromHex("#1a1a1a")).
		Padding(2).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#3366ff")).
		Width(m.width - 4)

	// Display content with instructions
	displayText := fmt.Sprintf(
		"[bold,cyan]Velvet Program Demo[]\n\n"+
			"Type something: [green]%s[]\n\n"+
			"[dim]Press ESC or Ctrl+C to quit[]",
		m.content,
	)

	rendered := container.Render(displayText)
	b.Place(rendered, buffer.AlignCenter, buffer.AlignMiddle)

	return b
}

func main() {
	m := model{
		content: "Hello, Velvet!",
	}

	p := velvet.NewProgram(m)

	if err := p.Run(); err != nil {
		panic(err)
	}
}
