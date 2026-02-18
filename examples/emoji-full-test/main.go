package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

func main() {
	fmt.Println("=== Testing Emoji Rendering ===\n")

	component := velvet.NewContainer().
		Foreground(color.FromHex("#fff")).
		Background(color.FromHex("#222")).
		Padding(1).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#00ff00")).
		Width(50)

	tests := []struct {
		name string
		text string
	}{
		{"Simple Emojis", "Hello 👋 World 🌍"},
		{"Multiple Emojis", "😀 😃 😄 😁 😆 😅 🤣"},
		{"Flags", "Flags: 🇺🇸 🇬🇧 🇯🇵 🇫🇷 🇩🇪"},
		{"Mixed Content", "ABC 123 👍 XYZ 789 🚀"},
		{"Animals", "🐶 🐱 🐭 🐹 🐰 🦊 🐻"},
		{"Food", "🍕 🍔 🍟 🌭 🍿 🧂"},
		{"Styled Text", "[bold]Bold 💪[] [cyan]Cyan 💙[] [red]Red ❤️[]"},
	}

	for _, test := range tests {
		fmt.Printf("Test: %s\n", test.name)
		result := component.Render(test.text)

		for y := 0; y < result.Height; y++ {
			for x := 0; x < result.Width; x++ {
				cell := result.Get(x, y)
				fmt.Print(cell.Content)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	fmt.Println("=== All Tests Complete ===")
}
