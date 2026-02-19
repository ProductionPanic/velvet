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

func (c CounterComponent) Render() velvet.Drawable {
	output := velvet.NewComponent().
		Foreground(color.FromHex("#fff")).
		Background(color.FromHex("#222")).
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

	output.Content(sb.String())
	return output
}
