package main

import (
	"fmt"
	"os"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/ansi"
	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
	"github.com/ProductionPanic/velvet/terminal"
)

func main() {
	w, h, err := terminal.GetSize()
	if err != nil {
		panic(err)
	}

	restore, err := terminal.RawMode()
	if err != nil {
		panic(err)
	}
	defer restore()

	ansi.Print(
		ansi.ClearScreen,
		ansi.EnterAltScreen,
		ansi.SetCursorPosition(1, 1),
		ansi.SetTitle("Velvet Demo"),
		ansi.HideCursor,
	)

	defer ansi.ResetAll()

	r := velvet.NewRenderer(os.Stdout, w, h)
	b := buffer.NewGrid(w, h)

	component := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Padding(1).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#3300ff")).
		Width(w / 4)

	b.Place(component.Render("Hello [red]velvet![] This is a demo of the Velvet [bold,cyan]TUI library[]. It supports text wrapping, borders, padding, and more!"), buffer.AlignCenter, buffer.AlignMiddle)

	r.Write(b)

	r.Flush()

	fmt.Scanln()
}
