package velvet

import (
	"strings"

	"github.com/ProductionPanic/velvet/style"
)

// View is a small wrapper around style.Style that provides helpers
// to render styled UI elements into a velvet.Buffer. We can't add
// methods directly to style.Style from here (different package), so
// this wrapper accepts a style.Style and produces a Buffer.
//
// The API is intentionally small: call NewView(s).Render(w,h,content)
// to get a *Buffer ready to be written by the renderer.
type View struct {
	style.Style
}

func NewView(s style.Style) View {
	return View{Style: s}
}

func (sv View) getSize(content string) (int, int) {
	//  if the size is set return it and else calculate it based on content and border
	width := sv.Width
	height := sv.Height

	if width <= 0 {
		lines := strings.Split(content, "\n")
		maxLineWidth := 0
		for _, line := range lines {
			if len(line) > maxLineWidth {
				maxLineWidth = len(line)
			}
		}
		width = maxLineWidth
	}

	if height <= 0 {
		lines := strings.Split(content, "\n")
		height = len(lines)
	}

	// add border size if needed
	if sv.Border.Top != "" || sv.Border.Bottom != "" {
		height += 2
	}
	if sv.Border.Left != "" || sv.Border.Right != "" {
		width += 2
	}

	return width, height
}

// Render creates a Buffer of size width x height containing the provided
// content laid out inside the area. If the Style has a border set (any of
// the border rune strings are non-empty), a one-character border is drawn
// around the content and the content is placed inside with a 1-character
// padding. Foreground and background colors from the Style are applied to
// all rendered cells via style.CellStyle.
func (sv View) Render(content string) *Buffer {
	width, height := sv.getSize(content)
	if width <= 0 || height <= 0 {
		return NewBuffer(0, 0)
	}

	buf := NewBuffer(width, height)

	// Build a cell style from the high-level Style colors
	cellStyle := style.CellStyle{
		Fg: sv.Fg,
		Bg: sv.Bg,
	}

	// Apply background/foreground to all cells first
	for i := range buf.Cells {
		buf.Cells[i].Style = cellStyle
		buf.Cells[i].Char = ' '
	}

	// Determine whether to draw a border
	hasBorder := sv.Border.Top != "" || sv.Border.Bottom != "" || sv.Border.Left != "" || sv.Border.Right != ""

	innerX, innerY := 0, 0
	innerW, innerH := width, height
	if hasBorder && width >= 2 && height >= 2 {
		// top and bottom
		for x := 0; x < width; x++ {
			if sv.Border.Top != "" {
				buf.Cells[buf.index(x, 0)] = Cell{Char: []rune(sv.Border.Top)[0], Style: cellStyle}
			}
			if sv.Border.Bottom != "" {
				buf.Cells[buf.index(x, height-1)] = Cell{Char: []rune(sv.Border.Bottom)[0], Style: cellStyle}
			}
		}
		// left and right
		for y := 1; y < height-1; y++ {
			if sv.Border.Left != "" {
				buf.Cells[buf.index(0, y)] = Cell{Char: []rune(sv.Border.Left)[0], Style: cellStyle}
			}
			if sv.Border.Right != "" {
				buf.Cells[buf.index(width-1, y)] = Cell{Char: []rune(sv.Border.Right)[0], Style: cellStyle}
			}
		}
		// corners
		if sv.Border.TopLeft != "" {
			buf.Cells[buf.index(0, 0)] = Cell{Char: []rune(sv.Border.TopLeft)[0], Style: cellStyle}
		}
		if sv.Border.TopRight != "" {
			buf.Cells[buf.index(width-1, 0)] = Cell{Char: []rune(sv.Border.TopRight)[0], Style: cellStyle}
		}
		if sv.Border.BottomLeft != "" {
			buf.Cells[buf.index(0, height-1)] = Cell{Char: []rune(sv.Border.BottomLeft)[0], Style: cellStyle}
		}
		if sv.Border.BottomRight != "" {
			buf.Cells[buf.index(width-1, height-1)] = Cell{Char: []rune(sv.Border.BottomRight)[0], Style: cellStyle}
		}

		innerX, innerY = 1, 1
		innerW, innerH = width-2, height-2
	}

	// Place the content into the inner area (wrap lines by width)
	lines := strings.Split(content, "\n")
	for y := 0; y < innerH && y < len(lines); y++ {
		line := lines[y]
		runes := []rune(line)
		for x := 0; x < innerW && x < len(runes); x++ {
			buf.Cells[buf.index(innerX+x, innerY+y)] = Cell{Char: runes[x], Style: cellStyle}
		}
	}

	return buf
}
