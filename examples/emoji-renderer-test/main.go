package main

import (
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
		ansi.SetTitle("Emoji Renderer Test"),
		ansi.HideCursor,
	)

	defer func() {
		ansi.ResetAll()
		ansi.Print(ansi.ShowCursor)
	}()

	r := velvet.NewRenderer(os.Stdout, w, h)
	b := buffer.NewGrid(w, h)

	component := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Padding(2).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#3300ff")).
		Width(w / 3)

	text := `[bold,cyan]Emoji Test 🎉[]

Hello 👋 World 🌍!
Emojis work: 😀✨
Flags: 🇺🇸 🇬🇧 🇯🇵
Thumbs up: 👍
Rocket: 🚀`

	b.Place(component.Render(text), buffer.AlignCenter, buffer.AlignMiddle)

	r.Write(b)
	r.Flush()

	// Wait for user input
	input := make([]byte, 1)
	os.Stdin.Read(input)
}
