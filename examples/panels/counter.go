package main

import (
	"fmt"
	"strings"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

type CounterComponent struct {
	Width, Height int
	Value         int
	IsActive      bool
	Title         string
}

func (c CounterComponent) Update(msg velvet.Msg) (CounterComponent, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.KeyMsg:
		switch msg.String() {
		case "up":
			c.Value++
		case "down":
			c.Value--
		case " ":
			c.Value = 0
		}
	}
	return c, nil
}

func (c CounterComponent) View() *velvet.Container {
	output := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Background(color.FromHex("#222")).
		Padding(1, 2).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#555")).
		Width(c.Width).
		Height(c.Height).
		BorderTextTop(c.Title)

	if c.IsActive {
		output.BorderColor(color.FromHex("#00ff00"))
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Counter: %d\n", c.Value))
	sb.WriteString("Use Up/Down to change the counter.\n")
	sb.WriteString("Press Space to reset.\n")
	sb.WriteString("Press Tab to switch panels.\n")
	sb.WriteString("Press Ctrl+C or Esc to quit.")

	output.Content(sb.String())
	return output
}
