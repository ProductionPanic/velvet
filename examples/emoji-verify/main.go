package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

func main() {
	component := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Background(color.FromHex("#000")).
		Padding(2).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#3300ff")).
		Width(60)

	tests := []string{
		"Hello 👋 World 🌍!",
		"Emojis: 😀🎉✨",
		"Mixed: ABC 👍 123 🚀",
		"Flags: 🇺🇸 🇬🇧 🇯🇵",
	}

	for _, test := range tests {
		fmt.Printf("\nTest: %s\n", test)
		result := component.Render(test)

		for y := 0; y < result.Height; y++ {
			for x := 0; x < result.Width; x++ {
				cell := result.Get(x, y)
				fmt.Print(cell.Content)
			}
			fmt.Println()
		}
	}
}
