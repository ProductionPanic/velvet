package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
)

type model struct {
	width, height int
	leftCounter   CounterComponent
	rightCounter  CounterComponent
	modal         ModalComponent
	showModal     bool
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
		m.leftCounter.Width, m.leftCounter.Height = m.width/2, m.height
		m.rightCounter.Width, m.rightCounter.Height = m.width/2, m.height
		return m, nil

	case velvet.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.isLeftActive {
				m.isLeftActive = false
				m.leftCounter.IsActive = false
				m.rightCounter.IsActive = true
			} else {
				m.isLeftActive = true
				m.leftCounter.IsActive = true
				m.rightCounter.IsActive = false
			}
			return m, nil
		case "enter":
			m.showModal = !m.showModal
			return m, nil
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	}

	var cmd velvet.Cmd
	if m.isLeftActive {
		m.leftCounter, cmd = m.leftCounter.Update(msg)
	} else {
		m.rightCounter, cmd = m.rightCounter.Update(msg)
	}

	return m, cmd
}

func (m model) View() *buffer.Grid {
	b := buffer.NewGrid(m.width, m.height)

	b.Place(m.leftCounter.View().Render(), buffer.AlignLeft, buffer.AlignTop)
	b.Place(m.rightCounter.View().Render(), buffer.AlignRight, buffer.AlignTop)

	if m.showModal {
		b.Place(m.modal.View().Render(), buffer.AlignCenter, buffer.AlignMiddle)
	}

	return b
}

func initialModel() model {
	return model{
		leftCounter:  CounterComponent{Title: "Left Counter", IsActive: true},
		rightCounter: CounterComponent{Title: "Right Counter"},
		modal: ModalComponent{
			Title:   "Modal Title",
			Content: "This is a modal. Press Enter to close.",
		},
		isLeftActive: true,
	}
}

func main() {
	p := velvet.NewProgram(initialModel())
	if err := p.Run(); err != nil {
		panic(err)
	}
}
