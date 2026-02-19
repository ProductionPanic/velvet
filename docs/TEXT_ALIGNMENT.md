# Text Alignment

The `TextAlign()` method allows you to control how text is horizontally aligned within a Component.

## Overview

By default, text in components is left-aligned. You can now set text alignment to:

- **`AlignLeft`** - Text starts from the left edge (default)
- **`AlignCenter`** - Text is centered horizontally
- **`AlignRight`** - Text is aligned to the right edge

## Usage

```go
import "github.com/ProductionPanic/velvet"

// Left-aligned (default)
component := velvet.NewComponent().
    Content("Left aligned text").
    TextAlign(velvet.AlignLeft)

// Center-aligned
component := velvet.NewComponent().
    Content("Center aligned text").
    TextAlign(velvet.AlignCenter)

// Right-aligned
component := velvet.NewComponent().
    Content("Right aligned text").
    TextAlign(velvet.AlignRight)
```

## Features

- Works with all **Component** instances (borders, padding, colors all supported)
- Affects **multiline text** - each line is aligned independently
- Respects **padding and borders** - alignment is calculated within the content area
- Part of the **fluent API** - chain with other styling methods

## Example

See `examples/textalign/main.go` for a complete working example that demonstrates all three alignment options side-by-side.

To run the example:
```bash
cd examples/textalign
go run main.go
```

Press `Esc` or `Ctrl+C` to exit.

## Implementation Details

The alignment is applied during text rendering in the `RenderWithContent()` method. The text offset is calculated based on:

1. The width of the current line
2. The available content width
3. The selected alignment mode
4. Padding and border widths

This ensures proper alignment regardless of component size or styling.

