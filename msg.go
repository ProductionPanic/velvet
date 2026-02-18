package velvet

import "time"

// Msg represents an event that can be sent to the program's Update function
type Msg interface{}

// KeyMsg represents a keyboard event
type KeyMsg struct {
	Type  KeyType
	Runes []rune
	Alt   bool
}

// KeyType represents a key press type
type KeyType int

const (
	KeyRunes KeyType = iota
	KeyEnter
	KeyBackspace
	KeyDelete
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyTab
	KeyShiftTab
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyEsc
	KeyCtrlC
	KeyCtrlD
	KeyCtrlZ
	KeySpace
)

// String returns a string representation of the key
func (k KeyMsg) String() string {
	if k.Type == KeyRunes {
		return string(k.Runes)
	}
	return ""
}

// WindowSizeMsg is sent when the terminal window is resized
type WindowSizeMsg struct {
	Width  int
	Height int
}

// QuitMsg is sent when the program should quit
type QuitMsg struct{}

// TickMsg is sent on a timer
type TickMsg struct {
	Time time.Time
}

// ErrMsg represents an error message
type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string {
	return e.Err.Error()
}
