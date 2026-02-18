package velvet

import (
	"io"
)

// readInput reads from the input and sends KeyMsg events
func readInput(msgs chan<- Msg, input io.Reader) {
	buf := make([]byte, 256)
	for {
		n, err := input.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			msgs <- ErrMsg{Err: err}
			return
		}

		if n == 0 {
			continue
		}

		msg := parseInput(buf[:n])
		if msg != nil {
			msgs <- msg
		}
	}
}

// parseInput converts raw bytes into KeyMsg
func parseInput(b []byte) Msg {
	if len(b) == 0 {
		return nil
	}

	// Single byte keys
	if len(b) == 1 {
		switch b[0] {
		case 3: // Ctrl+C
			return KeyMsg{Key: "ctrl+c"}
		case 4: // Ctrl+D
			return KeyMsg{Key: "ctrl+d"}
		case 9: // Tab
			return KeyMsg{Key: "tab"}
		case 13, 10: // Enter/Return
			return KeyMsg{Key: "enter"}
		case 27: // Escape
			return KeyMsg{Key: "esc"}
		case 127, 8: // Backspace/Delete
			return KeyMsg{Key: "backspace"}
		case 26: // Ctrl+Z
			return KeyMsg{Key: "ctrl+z"}
		case 32: // Space
			return KeyMsg{Key: " ", Runes: []rune{' '}}
		default:
			r := rune(b[0])
			return KeyMsg{Key: string(r), Runes: []rune{r}}
		}
	}

	// ANSI escape sequences
	if b[0] == 27 {
		if len(b) == 2 {
			// Alt + key
			r := rune(b[1])
			return KeyMsg{Key: string(r), Runes: []rune{r}, Alt: true}
		}

		if len(b) >= 3 && b[1] == '[' {
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
			case '3':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Key: "delete"}
				}
			case '5':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Key: "pgup"}
				}
			case '6':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Key: "pgdown"}
				}
			}
		}
	}

	// Multi-byte UTF-8 character
	runes := []rune(string(b))
	return KeyMsg{Key: string(runes), Runes: runes}
}
