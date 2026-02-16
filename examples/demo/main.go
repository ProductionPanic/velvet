package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/ansi"
	"github.com/ProductionPanic/velvet/buffer"
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

	done := make(chan bool)

	var colorSpeed float64 = 1.0
	var offset float64 = 0.0

	currentInterval := 50

	go func() {
		buf := make([]byte, 3)
		for {
			n, _ := os.Stdin.Read(buf)
			if n == 0 {
				continue
			}

			// Handle 'q' or Ctrl+C
			if buf[0] == 'q' || buf[0] == 3 {
				done <- true
				return
			}

			if buf[0] == 'r' || buf[0] == 'R' {
				colorSpeed *= -1 // Reverse direction
			}

			// Check for Escape Sequence (Arrows)
			if n == 3 && buf[0] == 27 && buf[1] == '[' {
				switch buf[2] {
				case 'A': // UP Arrow
					colorSpeed += 0.1
				case 'B': // DOWN Arrow
					colorSpeed -= 0.1
				}
			}
		}
	}()

	input := "Hello, Velvet animation! Use UP/DOWN arrows to change speed, 'R' to reverse, and 'Q' to quit."
	defaultStyle := velvet.NewCellStyle().SetBold(true)

	ticker := time.NewTicker(time.Millisecond * time.Duration(currentInterval))
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			offset += colorSpeed

			b.Clear()

			for i, r := range input {

				shift := int(float64(i) + offset)

				style := defaultStyle.Copy().SetFg(
					uint8((shift*20)%256), // R
					uint8((shift*40)%256), // G
					uint8((shift*60)%256), // B
				)

				startCellY := h/2 - 1
				sinValue := math.Sin(float64(startCellY)*offset + float64(i))
				cellY := startCellY + int(sinValue*float64(h/4)) // Vertical wave effect

				b.Set(i, cellY, buffer.NewCell(r, style))
			}
			// Display current speed for feedback
			speedMsg := fmt.Sprintf("Speed: %.1f", colorSpeed)
			for i, char := range speedMsg {
				b.Set(i+2, h-2, buffer.NewCell(char, velvet.NewCellStyle().SetFg(200, 200, 200)))
			}

			r.Write(b)
			r.Flush()
		}
	}

	fmt.Scanln()
}
