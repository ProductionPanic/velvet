package main

import (
	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

type ModalComponent struct {
	Title   string
	Content string
	IsError bool
}

func (m ModalComponent) Update(msg velvet.Msg) (ModalComponent, velvet.Cmd) {
	return m, nil
}

func (m ModalComponent) View() *velvet.Container {
	c := velvet.NewContainer().
		Border(style.RoundedBorder()).
		BorderTextTop(m.Title).
		BorderColor(color.FromHex("#00ff00"))

	if m.IsError {
		c.Background(color.FromHex("#ff0000"))
		c.Foreground(color.FromHex("#fff"))
	}

	c.Content(m.Content)

	return c
}
