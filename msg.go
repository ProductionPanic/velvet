package velvet

import "time"

// Msg represents an event that can be sent to the program's Update function
type Msg interface{}

// KeyMsg represents a keyboard event
type KeyMsg struct {
	Key   string
	Runes []rune
	Alt   bool
}

// String returns a string representation of the key
func (k KeyMsg) String() string {
	if k.Alt && k.Key != "" {
		return "alt+" + k.Key
	}
	return k.Key
}

// Matches checks if the key matches any of the provided key strings
func (k KeyMsg) Matches(keys ...string) bool {
	keyStr := k.String()
	for _, key := range keys {
		if keyStr == key {
			return true
		}
	}
	return false
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
