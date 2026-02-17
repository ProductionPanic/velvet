package style

type TextWrap int

const (
	BreakSpaces TextWrap = iota // BreakSpaces will break lines at spaces, and if a word is too long to fit on a line, it will be moved to the next line. if it doesnt fit then well break it
	BreakAny                    // dont care about spaces, just break at the end of the line
)
