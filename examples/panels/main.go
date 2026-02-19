package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
)

type model struct {
	width, height int
	counters      []CounterComponent
	activeCounter int
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
			if m.activeCounter >= len(m.counters) {
				m.activeCounter = 0
			}
			activeChanged = true
		case "-": // remove the current counter
			if len(m.counters) > 0 {
				m.counters = append(m.counters[:m.activeCounter], m.counters[m.activeCounter+1:]...)
				if m.activeCounter >= len(m.counters) {
					m.activeCounter = 0
				}
				activeChanged = true
			}
		case "+": // add a new counter
			newCounter := CounterComponent{
				Title: fmt.Sprintf("Counter %d", len(m.counters)+1),
			}
			m.counters = append(m.counters, newCounter)
			m.activeCounter = len(m.counters) - 1
			activeChanged = true
		case "shift+tab":
			m.activeCounter--
			if m.activeCounter < 0 {
				m.activeCounter = len(m.counters) - 1
			}
			activeChanged = true
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	}

	if activeChanged {
		for i := range m.counters {
			m.counters[i].IsActive = i == m.activeCounter
		}
		return m, nil
	}

	var cmd velvet.Cmd
	for i := range m.counters {
		if i == m.activeCounter {
			var cCmd velvet.Cmd
			m.counters[i], cCmd = m.counters[i].Update(msg)
			if cCmd != nil {
				cmd = cCmd
			}
		}
	}

	return m, cmd
}

func (m model) Render() velvet.Drawable {
	views := []velvet.Drawable{}
	for i := range m.counters {
		views = append(views, m.counters[i].Render())
	}
	return velvet.NewFlexContainer().
		SetSize(m.width, m.height).
		Add(views...)
}

func initialModel() model {
	counters := []CounterComponent{
		{Title: "Left Counter", IsActive: true},
		{Title: "Middle Counter"},
		{Title: "Right Counter"},
	}

	return model{
		counters:      counters,
		activeCounter: 0,
	}
}

func main() {
	p := velvet.NewProgram(initialModel())
	if err := p.Run(); err != nil {
		panic(err)
	}
}
