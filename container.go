package velvet

import (
	"strings"

	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
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

func (c *Container) getLongestLine(input [][]buffer.Cell) int {
	maxWidth := 0

	for _, line := range input {
		lineWidth := 0
		for _, cell := range line {
			lineWidth += cell.Width()
		}
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

func (c *Container) updateStyle(currentStyle *style.CellStyle, styleValue string) *style.CellStyle {
	parts := strings.Split(styleValue, ",") // Split the style value by comma to support multiple styles in one tag (e.g. [red,bold])

	return style.ApplyTemplateStyle(parts, currentStyle.Copy())
}

func (c *Container) parseLines(input string) [][]buffer.Cell {
	lines := strings.Split(input, "\n")
	var parsedLines [][]buffer.Cell

	// here we place the current style that is set within the input string (e.g. [red], [bold], [#f03], [bg#fff], etc.) and we update it whenever we encounter a new style tag in the input string.
	currentStyle := style.NewCellStyle()

	for _, line := range lines {
		var cells []buffer.Cell
		var inStyle bool
		var styleValue string
		var currentStyleValue string

		// Use grapheme cluster iterator to properly handle emojis and multi-codepoint characters
		gr := uniseg.NewGraphemes(line)
		for gr.Next() {
			cluster := gr.Str()

			// Check if this is a style tag opener
			if cluster == "[" {
				inStyle = true
				styleValue = ""
			} else if cluster == "]" && inStyle {
				inStyle = false
				currentStyleValue = styleValue
			} else if inStyle {
				styleValue += cluster
			} else {
				cells = append(cells, buffer.Cell{
					Content: cluster,
					Style:   c.updateStyle(currentStyle.Copy(), currentStyleValue),
				})
			}
		}
		parsedLines = append(parsedLines, cells)
	}

	return parsedLines
}

func (c *Container) applyTextWrap(lines [][]buffer.Cell, fixedContentWidth int, fixedContentHeight int) [][]buffer.Cell {
	var wrappedLines [][]buffer.Cell

	if fixedContentWidth > 0 && fixedContentHeight == 0 { // only width constraint
		if c.textWrap == style.BreakSpaces {
			lastSpaceIndex := -1
			lineWidth := 0
			currentLine := []buffer.Cell{}

			for _, line := range lines {
				for i, cell := range line {
					cellWidth := cell.Width()
					if cellWidth > 0 {
						if lineWidth+cellWidth > fixedContentWidth {
							if lastSpaceIndex != -1 {
								wrappedLines = append(wrappedLines, currentLine[:lastSpaceIndex])
								currentLine = currentLine[lastSpaceIndex+1:]
								lineWidth = 0
								for _, c := range currentLine {
									lineWidth += c.Width()
								}
							} else {
								wrappedLines = append(wrappedLines, currentLine)
								currentLine = []buffer.Cell{}
								lineWidth = 0
							}
						}
						currentLine = append(currentLine, cell)
						lineWidth += cellWidth

						if cell.Content == "" {
							lastSpaceIndex = len(currentLine) - 1
						}
						if i == len(line)-1 && len(currentLine) > 0 {
							wrappedLines = append(wrappedLines, currentLine)
							currentLine = []buffer.Cell{}
							lineWidth = 0
						}
					}
				}
				if len(currentLine) > 0 {
					wrappedLines = append(wrappedLines, currentLine)
					currentLine = []buffer.Cell{}
					lineWidth = 0
				}
			}
		} else if c.textWrap == style.BreakAny {
			lineWidth := 0
			currentLine := []buffer.Cell{}

			for _, line := range lines {
				for _, cell := range line {
					cellWidth := cell.Width()
					if cellWidth > 0 {
						if lineWidth+cellWidth > fixedContentWidth {
							wrappedLines = append(wrappedLines, currentLine)
							currentLine = []buffer.Cell{}
							lineWidth = 0
						}
						currentLine = append(currentLine, cell)
						lineWidth += cellWidth
					} else {
						currentLine = append(currentLine, cell)
					}
				}

				if len(currentLine) > 0 {
					wrappedLines = append(wrappedLines, currentLine)
					currentLine = []buffer.Cell{}
					lineWidth = 0
				}
			}
		}
	} else if fixedContentHeight > 0 && fixedContentWidth == 0 { // only height constraint
		if len(lines) > fixedContentHeight {
			wrappedLines = lines[:fixedContentHeight]

		} else {
			wrappedLines = lines
		}
	} else if fixedContentHeight > 0 && fixedContentWidth > 0 { // both width and height constraint
		// first we apply the width constraint and then we apply the height constraint on the wrapped lines.
		var tempLines [][]buffer.Cell

		if c.textWrap == style.BreakSpaces {
			lastSpaceIndex := -1
			lineWidth := 0
			currentLine := []buffer.Cell{}
			lastSpaceIndex = len(currentLine) - 1

			for _, line := range lines {
				for i, cell := range line {
					cellWidth := cell.Width()
					if cellWidth > 0 {
						if lineWidth+cellWidth > fixedContentWidth {
							if lastSpaceIndex != -1 {
								tempLines = append(tempLines, currentLine[:lastSpaceIndex])
								currentLine = currentLine[lastSpaceIndex+1:]
								lineWidth = 0
								for _, c := range currentLine {
									lineWidth += c.Width()
								}
							} else {
								tempLines = append(tempLines, currentLine)
								currentLine = []buffer.Cell{}
								lineWidth = 0
							}
						}
						currentLine = append(currentLine, cell)
						lineWidth += cellWidth

						if cell.Content == "" {
							lastSpaceIndex = len(currentLine) - 1
						}
						if i == len(line)-1 && len(currentLine) > 0 {
							tempLines = append(tempLines, currentLine)
							currentLine = []buffer.Cell{}
							lineWidth = 0
						}
					}
				}

				if len(currentLine) > 0 {
					tempLines = append(tempLines, currentLine)
					currentLine = []buffer.Cell{}
					lineWidth = 0
				}
			}

		} else if c.textWrap == style.BreakAny {
			lineWidth := 0
			currentLine := []buffer.Cell{}

			for _, line := range lines {
				for _, cell := range line {
					cellWidth := cell.Width()
					if cellWidth > 0 {
						if lineWidth+cellWidth > fixedContentWidth {
							tempLines = append(tempLines, currentLine)
							currentLine = []buffer.Cell{}
							lineWidth = 0
						}
						currentLine = append(currentLine, cell)
						lineWidth += cellWidth
					} else {
						currentLine = append(currentLine, cell)
					}
				}

				if len(currentLine) > 0 {
					tempLines = append(tempLines, currentLine)
					currentLine = []buffer.Cell{}
					lineWidth = 0
				}
			}
		}

		if len(tempLines) > fixedContentHeight {
			wrappedLines = tempLines[:fixedContentHeight]
		} else {
			wrappedLines = tempLines
		}
	} else {
		wrappedLines = lines
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

	lines := c.parseLines(input)
	parsedLines := c.applyTextWrap(lines, fixedContentWidth, fixedContentHeight)

	// now we have the parsed lines, we can create a grid and render the container with the parsed lines.
	contentWidth := c.getLongestLine(parsedLines)
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

	grid.SetEmptyCells("", &style.CellStyle{Fg: c.foreground, Bg: c.background})

	// set content
	for y, line := range parsedLines {
		col := 0 // track actual column position (accounts for wide characters)
		for _, r := range line {
			// Set the main cell
			grid.Set(col+c.getExtraWidth()/2, y+c.getExtraHeight()/2, buffer.Cell{
				Content: r.Content,
				Style:   r.Style,
			})

			// For wide characters (width > 1), fill continuation cells with spaces
			// to prevent overlap
			charWidth := r.Width()
			for i := 1; i < charWidth; i++ {
				grid.Set(col+i+c.getExtraWidth()/2, y+c.getExtraHeight()/2, buffer.Cell{
					Content: "",
					Style:   r.Style,
				})
			}

			col += charWidth // advance by the character's display width
		}
	}

	// set border
	if c.border != (style.BorderStyle{}) {
		for x := 0; x < gridWidth; x++ {
			grid.Set(x, 0, buffer.Cell{Content: c.border.Top, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(x, gridHeight-1, buffer.Cell{Content: c.border.Bottom, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
		for y := 0; y < gridHeight; y++ {
			grid.Set(0, y, buffer.Cell{Content: c.border.Left, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(gridWidth-1, y, buffer.Cell{Content: c.border.Right, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
		grid.Set(0, 0, buffer.Cell{Content: c.border.TopLeft, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(gridWidth-1, 0, buffer.Cell{Content: c.border.TopRight, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(0, gridHeight-1, buffer.Cell{Content: c.border.BottomLeft, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		grid.Set(gridWidth-1, gridHeight-1, buffer.Cell{Content: c.border.BottomRight, Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
	}

	return grid
}
