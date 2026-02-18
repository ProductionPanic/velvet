package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

func main() {
	spaceWidth := runewidth.RuneWidth(' ')
	fmt.Println(spaceWidth)
}
