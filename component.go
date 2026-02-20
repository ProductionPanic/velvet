package velvet

import (
	"strings"

	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// TextAlignment defines how text should be aligned horizontally within a component.
type TextAlignment int

const (
	// AlignLeft aligns text to the left (default).
	AlignLeft TextAlignment = iota
	// AlignCenter centers the text.
	AlignCenter
	// AlignRight aligns text to the right.
	AlignRight
)

type Component struct {
	padding struct {
		top, right, bottom, left int
	}
	border           style.BorderStyle
	content          string
	borderTextTop    string // for when you want to set an inline text in the border (e.g. ──[ Title ]──)
	borderTextBott   string // same but for in the bottom
	borderTextMargin int    // margin between border text and border line
	borderColor      color.Color
	foreground       color.Color
	background       color.Color
	textWrap         style.TextWrap
	textAlign        TextAlignment
	width            int
	height           int
	zoneID           string
}

func NewComponent() *Component {
	return &Component{
		textWrap:         style.BreakSpaces,
		borderTextMargin: 1,
	}
}

func (c *Component) Copy() *Component {
	newContainer := *c
	return &newContainer
}

func (c *Component) Content(content string) *Component {
	c.content = content
	return c
}

func (c *Component) SetZoneID(id string) *Component {
	c.zoneID = id
	return c
}

func (c *Component) Padding(p ...int) *Component {
	if len(p) == 1 {
		padding := p[0]
		c.padding.top = padding
		c.padding.right = padding
		c.padding.bottom = padding
		c.padding.left = padding
	} else if len(p) == 2 {
		c.padding.top = p[0]
		c.padding.bottom = p[0]
		c.padding.right = p[1]
		c.padding.left = p[1]
	} else if len(p) == 4 {
		c.padding.top = p[0]
		c.padding.right = p[1]
		c.padding.bottom = p[2]
		c.padding.left = p[3]
	}
	return c
}

func (c *Component) BorderTextTop(text string) *Component {
	c.borderTextTop = text
	return c
}

func (c *Component) BorderTextBottom(text string) *Component {
	c.borderTextBott = text
	return c
}

func (c *Component) PaddingTop(p int) *Component {
	c.padding.top = p
	return c
}

func (c *Component) PaddingRight(p int) *Component {
	c.padding.right = p
	return c
}

func (c *Component) PaddingBottom(p int) *Component {
	c.padding.bottom = p
	return c
}

func (c *Component) PaddingLeft(p int) *Component {
	c.padding.left = p
	return c
}

func (c *Component) Border(b style.BorderStyle) *Component {
	c.border = b
	return c
}

func (c *Component) BorderColor(col color.Color) *Component {
	c.borderColor = col
	return c
}

func (c *Component) Foreground(col color.Color) *Component {
	c.foreground = col
	return c
}

func (c *Component) Background(col color.Color) *Component {
	c.background = col
	return c
}

func (c *Component) Size(width, height int) *Component {
	c.width = width
	c.height = height
	return c
}

func (c *Component) Width(width int) *Component {
	c.width = width
	return c
}

func (c *Component) Height(height int) *Component {
	c.height = height
	return c
}

func (c *Component) TextAlign(align TextAlignment) *Component {
	c.textAlign = align
	return c
}

func (c *Component) getTextXOffset(lineWidth, contentWidth int) int {
	switch c.textAlign {
	case AlignLeft:
		return c.padding.left
	case AlignCenter:
		return (contentWidth-lineWidth)/2 + c.padding.left
	case AlignRight:
		return contentWidth - lineWidth + c.padding.left
	default:
		return c.padding.left
	}
}

func (c *Component) getLongestLine(input [][]buffer.Cell) int {
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

func (c *Component) getExtraWidth() int {
	paddingWidth := c.padding.left + c.padding.right
	borderWidth := 0

	if c.border != (style.BorderStyle{}) {
		borderWidth = 2 // Assuming a simple border that adds 1 character on each side
	}

	return paddingWidth + borderWidth
}

func (c *Component) getExtraHeight() int {
	paddingHeight := c.padding.top + c.padding.bottom
	borderHeight := 0

	if c.border != (style.BorderStyle{}) {
		borderHeight = 2 // Assuming a simple border that adds 1 character on top and bottom
	}

	return paddingHeight + borderHeight
}

func (c *Component) getBorderWidth() int {
	if c.border != (style.BorderStyle{}) {
		return 1 // Border adds 1 character on the left side
	}
	return 0
}

func (c *Component) wrapText(text string, maxWidth int) []string {
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

func (c *Component) updateStyle(currentStyle *style.CellStyle, styleValue string) *style.CellStyle {
	parts := strings.Split(styleValue, ",") // Split the style value by comma to support multiple styles in one tag (e.g. [red,bold])

	return style.ApplyTemplateStyle(parts, currentStyle.Copy())
}

func (c *Component) getContainerCellStyle() *style.CellStyle {
	// build a cell style based on the container's foreground and background colors
	return &style.CellStyle{
		Fg: c.foreground,
		Bg: c.background,
	}
}

func (c *Component) parseLines(input string) [][]buffer.Cell {
	lines := strings.Split(input, "\n")
	var parsedLines [][]buffer.Cell

	// here we place the current style that is set within the input string (e.g. [red], [bold], [#f03], [bg#fff], etc.) and we update it whenever we encounter a new style tag in the input string.
	currentStyle := c.getContainerCellStyle()

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

func (c *Component) applyTextWrap(lines [][]buffer.Cell, fixedContentWidth int, fixedContentHeight int) [][]buffer.Cell {
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

						if cell.Content == " " {
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

						if cell.Content == " " {
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

// Render implements the Drawable interface
func (c *Component) Render() *buffer.Grid {
	return c.RenderWithContent().WithZoneID(c.zoneID)
}

// RenderWithContent renders the component with optional additional content
func (c *Component) RenderWithContent(s ...string) *buffer.Grid {
	input := c.content + strings.Join(s, "\n")
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

	// If fixed width is set, use it; otherwise use content width
	if fixedContentWidth > 0 {
		contentWidth = fixedContentWidth
	}

	contentHeight := len(parsedLines)

	// If fixed height is set, use it; otherwise use content height
	if fixedContentHeight > 0 {
		contentHeight = fixedContentHeight
	}

	gridWidth := contentWidth + c.getExtraWidth()
	gridHeight := contentHeight + c.getExtraHeight()

	grid := buffer.NewGrid(gridWidth, gridHeight)

	grid.SetEmptyCells(" ", &style.CellStyle{Fg: c.foreground, Bg: c.background})

	// set content
	for y, line := range parsedLines {
		// Calculate the x offset for this line based on alignment
		lineWidth := 0
		for _, r := range line {
			lineWidth += r.Width()
		}
		xOffset := c.getTextXOffset(lineWidth, contentWidth) + c.getBorderWidth()

		col := 0 // track actual column position (accounts for wide characters)
		for _, r := range line {
			// Set the main cell
			grid.Set(col+xOffset, y+c.getExtraHeight()/2, buffer.Cell{
				Content: r.Content,
				Style:   r.Style,
			})

			// For wide characters (width > 1), fill continuation cells with spaces
			// to prevent overlap
			charWidth := r.Width()
			for i := 1; i < charWidth; i++ {
				grid.Set(col+i+xOffset, y+c.getExtraHeight()/2, buffer.Cell{
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

	// if needed we set the border text (e.g. ──[ Title ]──)
	if c.borderTextTop != "" {
		borderText := c.borderTextTop
		borderTextWidth := runewidth.StringWidth(borderText)
		startX := (gridWidth-borderTextWidth)/2 - c.borderTextMargin
		for i, r := range borderText {
			grid.Set(startX+i, 0, buffer.Cell{Content: string(r), Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}

		// Fill spaces around the border text to prevent overlap with border lines
		for i := 0; i < c.borderTextMargin; i++ {
			grid.Set(startX-1-i, 0, buffer.Cell{Content: " ", Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(startX+borderTextWidth+i, 0, buffer.Cell{Content: " ", Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
	}

	if c.borderTextBott != "" {
		borderText := c.borderTextBott
		borderTextWidth := runewidth.StringWidth(borderText)
		startX := (gridWidth-borderTextWidth)/2 - c.borderTextMargin
		for i, r := range borderText {
			grid.Set(startX+i, gridHeight-1, buffer.Cell{Content: string(r), Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}

		// Fill spaces around the border text to prevent overlap with border lines
		for i := 0; i < c.borderTextMargin; i++ {
			grid.Set(startX-1-i, gridHeight-1, buffer.Cell{Content: " ", Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
			grid.Set(startX+borderTextWidth+i, gridHeight-1, buffer.Cell{Content: " ", Style: &style.CellStyle{Fg: c.borderColor, Bg: c.background}})
		}
	}

	return grid
}
