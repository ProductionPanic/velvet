package main

import (
	"fmt"
	"time"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style/color"
)

type model struct {
	width, height int
	tick          int
}

func (m model) Init() velvet.Cmd {
	return velvet.Tick(time.Second)
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case velvet.TickMsg:
		m.tick++
		return m, velvet.Tick(time.Second)

	case velvet.KeyMsg:
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

	// Fill with debug info
	msg := fmt.Sprintf("Size: %dx%d | Tick: %d | Press ESC to quit", m.width, m.height, m.tick)

	// Put text at top-left
	x := 2
	y := 2
	for i, ch := range msg {
		if x+i < m.width {
			b.SetRune(x+i, y, string(ch))
			cell := b.Get(x+i, y)
			cell.Style.Fg = color.FromHex("#00ff00")
			b.Set(x+i, y, cell)
		}
	}

	return b
}

func main() {
	p := velvet.NewProgram(model{})
	if err := p.Run(); err != nil {
		panic(err)
	}
}
