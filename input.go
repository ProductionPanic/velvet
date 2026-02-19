package velvet

import (
	"bufio"
	"fmt"
	"io"
)

// readInput now uses a buffered reader to handle multi-byte escape sequences
func readInput(msgs chan<- Msg, input io.Reader) {
	reader := bufio.NewReader(input)

	for {
		// Read the first byte
		b, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				msgs <- ErrMsg{Err: err}
			}
			return
		}

		// Handle Escape Sequences (Keys and Mouse)
		if b == 27 { // ESC
			// Peek to see if there are more bytes following immediately
			if reader.Buffered() > 0 {
				next, _ := reader.Peek(1)
				if next[0] == '[' {
					// This is a CSI sequence (Control Sequence Introducer)
					// We read until the "terminator" character (letters like M, m, A, B, etc.)
					seq := []byte{b}
					for {
						nb, _ := reader.ReadByte()
						seq = append(seq, nb)
						// SGR Mouse ends in 'M' or 'm'. Standard keys end in 'A'-'Z' or '~'
						if (nb >= 'A' && nb <= 'Z') || (nb >= 'a' && nb <= 'z') || nb == '~' {
							break
						}
					}
					msgs <- parseInput(seq)
					continue
				}

				// Handle Alt+Key (Esc followed by a single byte)
				nb, _ := reader.ReadByte()
				msgs <- KeyMsg{Key: string(nb), Runes: []rune{rune(nb)}, Alt: true}
				continue
			}
			// Just a plain Escape key
			msgs <- KeyMsg{Key: "esc"}
			continue
		}

		// Handle standard single-byte or UTF-8 multi-byte characters
		if b < 128 {
			msgs <- parseInput([]byte{b})
		} else {
			// Multi-byte UTF-8
			unquoted := []byte{b}
			// Basic logic to grab the rest of the UTF-8 sequence based on the leading byte
			var extra int
			if b >= 0xE0 {
				extra = 2
			} else if b >= 0xC0 {
				extra = 1
			}
			for i := 0; i < extra; i++ {
				nb, _ := reader.ReadByte()
				unquoted = append(unquoted, nb)
			}
			msgs <- KeyMsg{Key: string(unquoted), Runes: []rune(string(unquoted))}
		}
	}
}

// parseInput handles the logic for turning bytes into specific Msg types
func parseInput(b []byte) Msg {
	if len(b) == 0 {
		return nil
	}

	// Handle SGR Mouse Sequences: ^[[<button;x;yM or ^[[<button;x;ym
	if len(b) >= 6 && string(b[0:3]) == "\x1b[<" {
		var button, x, y int
		lastChar := b[len(b)-1]

		// Parse the numbers between '<' and the final 'M/m'
		_, err := fmt.Sscanf(string(b[3:len(b)-1]), "%d;%d;%d", &button, &x, &y)
		if err == nil {
			msg := MouseMsg{
				X:      x,
				Y:      y,
				Button: button,
			}

			// Button bits: 32 = Motion, 64 = Wheel Up, 65 = Wheel Down
			// Terminator: 'm' = Release, 'M' = Press/Motion
			if lastChar == 'm' {
				msg.Type = MouseRelease
			} else if button&32 != 0 {
				msg.Type = MouseMotion
			} else {
				msg.Type = MousePress
			}
			return msg
		}
	}

	// Single Byte / Control Keys
	if len(b) == 1 {
		switch b[0] {
		case 3:
			return KeyMsg{Key: "ctrl+c"}
		case 4:
			return KeyMsg{Key: "ctrl+d"}
		case 9:
			return KeyMsg{Key: "tab"}
		case 13, 10:
			return KeyMsg{Key: "enter"}
		case 127, 8:
			return KeyMsg{Key: "backspace"}
		case 32:
			return KeyMsg{Key: " ", Runes: []rune{' '}}
		default:
			r := rune(b[0])
			return KeyMsg{Key: string(r), Runes: []rune{r}}
		}
	}

	// Standard ANSI Cursor/Function Keys
	if b[0] == 27 && b[1] == '[' {
		switch b[2] {
		case 'A':
			return KeyMsg{Key: "up"}
		case 'B':
			return KeyMsg{Key: "down"}
		case 'C':
			return KeyMsg{Key: "right"}
		case 'D':
			return KeyMsg{Key: "left"}
		case 'H':
			return KeyMsg{Key: "home"}
		case 'F':
			return KeyMsg{Key: "end"}
		case 'Z':
			return KeyMsg{Key: "shift+tab"}
		}
	}

	return nil
}
