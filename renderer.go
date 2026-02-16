package velvet

import (
	"fmt"
	"io"
	"strings"

	"github.com/ProductionPanic/velvet/ansi"
)

type renderer struct {
	frontBuffer *Buffer
	backBuffer  *Buffer
	Out         io.Writer
	Width       int
	Height      int
}

func (r *renderer) Flush() {
	var sb strings.Builder
	changed := false
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			idx := y*r.Width + x
			newCell := r.backBuffer.Cells[idx]
			oldCell := r.frontBuffer.Cells[idx]

			if newCell != oldCell {
				changed = true
				// set cursor position
				sb.WriteString(ansi.SetCursorPosition(x, y))
				// apply style
				sb.WriteString(ansi.RenderCellStyle(newCell.Style))
				// write character
				sb.WriteRune(newCell.Char)
				// reset style
				sb.WriteString(ansi.Reset)

				// update front buffer
				r.frontBuffer.Cells[idx] = newCell
			}
		}
	}
	if !changed {
		return
	}

	fmt.Fprint(r.Out, sb.String())

	r.backBuffer.Clear()
}

func (r *renderer) Write(buf *Buffer) {
	// size is not the same
	if buf.Width != r.Width || buf.Height != r.Height {
		// crop the buffer to fit the renderer
		buf = buf.Crop(r.Width, r.Height)
	}

	r.backBuffer = buf
}

func newRenderer(out io.Writer, width, height int) renderer {
	return renderer{
		frontBuffer: NewBuffer(width, height),
		backBuffer:  NewBuffer(width, height),
		Out:         out,
		Width:       width,
		Height:      height,
	}
}
