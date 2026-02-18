package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
)

type model struct {
	width, height int
	leftCounter   CounterComponent
	middleCounter CounterComponent
	rightCounter  CounterComponent
	activeCounter int
	modal         ModalComponent
	showModal     bool
}

func (m model) Init() velvet.Cmd {
	return nil
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	var activeChanged bool = false

	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case velvet.KeyMsg:
		switch msg.String() {
		case "tab":
			m.activeCounter++
			if m.activeCounter > 2 {
				m.activeCounter = 0
			}
			activeChanged = true
		case "shift+tab":
			m.activeCounter--
			if m.activeCounter < 0 {
				m.activeCounter = 2
			}
			activeChanged = true
		case "enter":
			m.showModal = !m.showModal
			return m, nil
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	}

	if activeChanged {
		m.leftCounter.IsActive = m.activeCounter == 0
		m.middleCounter.IsActive = m.activeCounter == 1
		m.rightCounter.IsActive = m.activeCounter == 2
		return m, nil
	}

	var cmd velvet.Cmd
	switch m.activeCounter {
	case 0:
		m.leftCounter, cmd = m.leftCounter.Update(msg)
	case 1:
		m.middleCounter, cmd = m.middleCounter.Update(msg)
	case 2:
		m.rightCounter, cmd = m.rightCounter.Update(msg)
	}

	return m, cmd
}

func (m model) View() *buffer.Grid {
	b := velvet.NewFlexContainer().
		SetSize(m.width, m.height).
		Add(
			m.leftCounter.View(),
			m.middleCounter.View(),
			m.rightCounter.View(),
		).
		Render()

	if m.showModal {
		b.Place(m.modal.View().Render(), buffer.AlignCenter, buffer.AlignMiddle)
	}

	return b
}

func initialModel() model {
	return model{
		leftCounter:   CounterComponent{Title: "Left Counter", IsActive: true},
		middleCounter: CounterComponent{Title: "Middle Counter"},
		rightCounter:  CounterComponent{Title: "Right Counter"},
		modal: ModalComponent{
			Title:   "Modal Title",
			Content: "This is a modal. Press Enter to close.",
			IsError: false,
		},
	}
}

func main() {
	p := velvet.NewProgram(initialModel())
	if err := p.Run(); err != nil {
		panic(err)
	}
}
