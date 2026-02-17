package velvet

import (
	"strings"

	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
	"github.com/mattn/go-runewidth"
)

type Container struct {
	padding struct {
		top, right, bottom, left int
	}
	border      style.BorderStyle
	borderColor color.Color
	foreground  color.Color
	background  color.Color
	textWrap    style.TextWrap
	width       int
	height      int
}

func NewContainer() *Container {
	return &Container{
		textWrap: style.BreakSpaces,
	}
}

func (c *Container) Padding(p int) *Container {
	c.padding.top = p
	c.padding.right = p
	c.padding.bottom = p
	c.padding.left = p
	return c
}

func (c *Container) PaddingTop(p int) *Container {
	c.padding.top = p
	return c
}

func (c *Container) PaddingRight(p int) *Container {
	c.padding.right = p
	return c
}

func (c *Container) PaddingBottom(p int) *Container {
	c.padding.bottom = p
	return c
}

func (c *Container) PaddingLeft(p int) *Container {
	c.padding.left = p
	return c
}

func (c *Container) Border(b style.BorderStyle) *Container {
	c.border = b
	return c
}

func (c *Container) BorderColor(col color.Color) *Container {
	c.borderColor = col
	return c
}

func (c *Container) Foreground(col color.Color) *Container {
	c.foreground = col
	return c
}

func (c *Container) Background(col color.Color) *Container {
	c.background = col
	return c
}

func (c *Container) Size(width, height int) *Container {
	c.width = width
	c.height = height
	return c
}

func (c *Container) Width(width int) *Container {
	c.width = width
	return c
}

func (c *Container) Height(height int) *Container {
	c.height = height
	return c
}

func (c *Container) getLongestLine(input string) int {
	maxWidth := 0
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		lineWidth := runewidth.StringWidth(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return maxWidth
}

func (c *Container) getExtraWidth() int {
	paddingWidth := c.padding.left + c.padding.right
	borderWidth := 0

	if c.border != (style.BorderStyle{}) {
		borderWidth = 2 // Assuming a simple border that adds 1 character on each side
	}

	return paddingWidth + borderWidth
}

func (c *Container) getExtraHeight() int {
	paddingHeight := c.padding.top + c.padding.bottom
	borderHeight := 0

	if c.border != (style.BorderStyle{}) {
		borderHeight = 2 // Assuming a simple border that adds 1 character on top and bottom
	}

	return paddingHeight + borderHeight
}

func (c *Container) wrapText(text string, maxWidth int) []string {
	var wrappedLines []string
	currentLine := ""

	if c.textWrap == style.BreakSpaces {
		words := strings.Fields(text)
		for _, word := range words {
			if runewidth.StringWidth(currentLine)+runewidth.StringWidth(word)+1 > maxWidth {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = word
			} else {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			}
		}
		if currentLine != "" {
			wrappedLines = append(wrappedLines, currentLine)
		}

	} else if c.textWrap == style.BreakAny {
		runes := []rune(text)
		currentWidth := 0

		for _, r := range runes {
			runeWidth := runewidth.RuneWidth(r)
			if currentWidth+runeWidth > maxWidth {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = string(r)
				currentWidth = runeWidth
			} else {
				currentLine += string(r)
				currentWidth += runeWidth
			}
		}
		if currentLine != "" {
			wrappedLines = append(wrappedLines, currentLine)
		}
	}

	return wrappedLines
}

func (c *Container) Render(input string) *buffer.Grid {
	// first we check if there is any constraints on the width and height.
	var fixedContentWidth, fixedContentHeight int
	if c.width > 0 {
		fixedContentWidth = c.width - c.getExtraWidth()
	}

	if c.height > 0 {
		fixedContentHeight = c.height - c.getExtraHeight()
	}

	lines := strings.Split(input, "\n")
	var parsedLines []string

	if fixedContentWidth > 0 && fixedContentHeight == 0 { // only width constraint
		for _, line := range lines {
			if runewidth.StringWidth(line) > fixedContentWidth {
				parts := c.wrapText(line, fixedContentWidth) // wrap the line into multiple lines based on the fixed content width and the text wrap mode

				for _, part := range parts {
					parsedLines = append(parsedLines, part)
				}
			} else {
				parsedLines = append(parsedLines, line)
			}
		}
	} else if fixedContentHeight > 0 && fixedContentWidth == 0 { // only height constraint
		if len(lines) > fixedContentHeight {
			parsedLines = lines[:fixedContentHeight] // truncate the lines to fit the fixed content height
		}
	} else if fixedContentHeight > 0 && fixedContentWidth > 0 { // both width and height constraint
		for _, line := range lines {
			if runewidth.StringWidth(line) > fixedContentWidth {
				parts := c.wrapText(line, fixedContentWidth) // wrap the line into multiple lines based on the fixed content width and the text wrap mode

				for _, part := range parts {
					parsedLines = append(parsedLines, part)
				}

				if len(parsedLines) >= fixedContentHeight {
					parsedLines = parsedLines[:fixedContentHeight] // truncate the lines to fit the fixed content height
					break
				}
			} else {
				parsedLines = append(parsedLines, line)

				if len(parsedLines) >= fixedContentHeight {
					parsedLines = parsedLines[:fixedContentHeight] // truncate the lines to fit the fixed content height
					break
				}
			}
		}
	} else {
		parsedLines = lines
	}

	// now we have the parsed lines, we can create a grid and render the container with the parsed lines.
	contentWidth := c.getLongestLine(
		strings.Join(parsedLines, "\n"),
	)
	if fixedContentWidth > 0 && contentWidth > fixedContentWidth {
		contentWidth = fixedContentWidth
	}
	contentHeight := len(parsedLines)
	if fixedContentHeight > 0 && contentHeight > fixedContentHeight {
		contentHeight = fixedContentHeight
	}

	gridWidth := contentWidth + c.getExtraWidth()
	gridHeight := contentHeight + c.getExtraHeight()

	grid := buffer.NewGrid(gridWidth, gridHeight)

	grid.SetEmptyCells(' ', &style.CellStyle{Fg: c.foreground, Bg: c.background})

	// set content
	for y, line := range parsedLines {
		for x, r := range line {
			grid.Set(x+c.getExtraWidth()/2, y+c.getExtraHeight()/2, buffer.Cell{
				Rune: r,
				Style: &style.CellStyle{
					Fg: c.foreground,
					Bg: c.background,
				},
			})
		}
	}

	// set border
	if c.border != (style.BorderStyle{}) {
		for x := 0; x < gridWidth; x++ {
			grid.Set(x, 0, buffer.Cell{Rune: rune(c.border.Top), Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(x, gridHeight-1, buffer.Cell{Rune: c.border.Bottom, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
		for y := 0; y < gridHeight; y++ {
			grid.Set(0, y, buffer.Cell{Rune: c.border.Left, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(gridWidth-1, y, buffer.Cell{Rune: c.border.Right, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
		grid.Set(0, 0, buffer.Cell{Rune: c.border.TopLeft, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(gridWidth-1, 0, buffer.Cell{Rune: c.border.TopRight, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(0, gridHeight-1, buffer.Cell{Rune: c.border.BottomLeft, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(gridWidth-1, gridHeight-1, buffer.Cell{Rune: c.border.BottomRight, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
	}

	return grid
}
