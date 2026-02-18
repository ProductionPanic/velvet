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
			return KeyMsg{Type: KeyCtrlC}
		case 4: // Ctrl+D
			return KeyMsg{Type: KeyCtrlD}
		case 9: // Tab
			return KeyMsg{Type: KeyTab}
		case 13, 10: // Enter/Return
			return KeyMsg{Type: KeyEnter}
		case 27: // Escape
			return KeyMsg{Type: KeyEsc}
		case 127, 8: // Backspace/Delete
			return KeyMsg{Type: KeyBackspace}
		case 26: // Ctrl+Z
			return KeyMsg{Type: KeyCtrlZ}
		case 32: // Space
			return KeyMsg{Type: KeySpace, Runes: []rune{' '}}
		default:
			return KeyMsg{Type: KeyRunes, Runes: []rune{rune(b[0])}}
		}
	}

	// ANSI escape sequences
	if b[0] == 27 {
		if len(b) == 2 {
			// Alt + key
			return KeyMsg{Type: KeyRunes, Runes: []rune{rune(b[1])}, Alt: true}
		}

		if len(b) >= 3 && b[1] == '[' {
			switch b[2] {
			case 'A':
				return KeyMsg{Type: KeyUp}
			case 'B':
				return KeyMsg{Type: KeyDown}
			case 'C':
				return KeyMsg{Type: KeyRight}
			case 'D':
				return KeyMsg{Type: KeyLeft}
			case 'H':
				return KeyMsg{Type: KeyHome}
			case 'F':
				return KeyMsg{Type: KeyEnd}
			case 'Z':
				return KeyMsg{Type: KeyShiftTab}
			case '3':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Type: KeyDelete}
				}
			case '5':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Type: KeyPgUp}
				}
			case '6':
				if len(b) >= 4 && b[3] == '~' {
					return KeyMsg{Type: KeyPgDown}
				}
			}
		}
	}

	// Multi-byte UTF-8 character
	runes := []rune(string(b))
	return KeyMsg{Type: KeyRunes, Runes: runes}
}
