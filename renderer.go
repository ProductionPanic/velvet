package velvet

import (
	"fmt"
	"io"
	"strings"

	"github.com/ProductionPanic/velvet/ansi"
	"github.com/ProductionPanic/velvet/buffer"
)

type Renderer struct {
	Out           io.Writer
	front         *buffer.Grid
	back          *buffer.Grid
	width, height int
}

// replace the backbuffer
func (r *Renderer) Write(content *buffer.Grid) {
	r.back = content.Resize(r.width, r.height) // size to fit the renderer
}

func (r *Renderer) Clear() {
	r.back = buffer.NewGrid(r.width, r.height)
}

func (r *Renderer) outputCell(cell buffer.Cell, x, y int, sb *strings.Builder) {
	sb.WriteString(ansi.Combine(
		ansi.SetCursorPosition(x+1, y+1), // ANSI escape codes are 1-indexed)
		ansi.ForegroundColor(cell.Style.Fg),
		ansi.BackgroundColor(cell.Style.Bg),
	))

	if cell.Style.Bold {
		ansi.Print(ansi.Bold)
	}

	if cell.Style.Italic {
		ansi.Print(ansi.Italic)
	}

	if cell.Style.Underline {
		ansi.Print(ansi.Underline)
	}

	sb.WriteString(cell.Content)

	sb.WriteString(ansi.Reset) // Reset styles after each cell to avoid style bleed

}

func (r *Renderer) Flush() {
	var sb strings.Builder

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			frontCell := r.front.Get(x, y)
			backCell := r.back.Get(x, y)

			if frontCell != backCell {
				if backCell.Content != "" {
					r.outputCell(backCell, x, y, &sb)
				}
				r.front.Set(x, y, backCell)
			}
		}
	}

	fmt.Fprint(r.Out, sb.String())
}

func NewRenderer(out io.Writer, width, height int) *Renderer {
	front := buffer.NewGrid(width, height)
	back := buffer.NewGrid(width, height)

	return &Renderer{
		Out:    out,
		front:  front,
		back:   back,
		width:  width,
		height: height,
	}
}
