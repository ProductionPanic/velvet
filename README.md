# Velvet 🎭

A modern, buffer-based TUI (Terminal User Interface) library for Go with a bubbletea-inspired architecture.

## Features

✨ **Buffer-Based Rendering** - Work with cells and grids instead of strings for precise control  
🎨 **Rich Styling** - Colors, borders, padding, text wrapping, and template syntax  
⚡ **Event-Driven** - Message-based architecture for reactive applications  
⌨️ **Full Keyboard Support** - Arrow keys, Ctrl combinations, UTF-8 input  
📐 **Flexible Layout** - Grid system with alignment controls  
🔄 **Automatic Management** - Terminal setup, cleanup, and resize handling  

## Quick Start

```go
package main

import (
    "github.com/ProductionPanic/velvet"
    "github.com/ProductionPanic/velvet/buffer"
    "github.com/ProductionPanic/velvet/style"
    "github.com/ProductionPanic/velvet/style/color"
)

type model struct {
    width, height int
}

func (m model) Init() velvet.Cmd {
    return nil
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
    switch msg := msg.(type) {
    case velvet.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
    case velvet.KeyMsg:
        if msg.Type == velvet.KeyCtrlC || msg.Type == velvet.KeyEsc {
            return m, func() velvet.Msg { return velvet.Quit() }
        }
    }
    return m, nil
}

func (m model) View() *buffer.Grid {
    b := buffer.NewGrid(m.width, m.height)
    
    container := velvet.NewContainer().
        Foreground(color.FromHex("#ffffff")).
        Padding(2).
        Border(style.RoundedBorder()).
        BorderColor(color.FromHex("#3366ff")).
        Width(m.width / 2)
    
    text := "[bold,cyan]Hello, Velvet![]\\n\\nPress ESC to quit"
    b.Place(container.Render(text), buffer.AlignCenter, buffer.AlignMiddle)
    
    return b
}

func main() {
    p := velvet.NewProgram(model{})
    if err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Installation

```bash
go get github.com/ProductionPanic/velvet
```

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started quickly with examples
- **[PROGRAM.md](PROGRAM.md)** - Complete Program API documentation
- **[IMPLEMENTATION.md](IMPLEMENTATION.md)** - Implementation details and architecture

## Core Components

### Program API

The `Program` struct manages your application lifecycle:

- **Model Interface**: `Init()`, `Update(Msg)`, `View()`
- **Message System**: Keyboard, resize, timer, and custom messages
- **Command System**: Async operations with `Batch`, `Sequence`, `Tick`
- **Automatic Setup**: Raw mode, alt screen, cursor management

### Rendering System

- **Grid**: 2D cell buffer with width/height
- **Cell**: Individual character with style (color, bold, italic, underline)
- **Container**: Component with borders, padding, alignment, text wrapping
- **Renderer**: Efficient diff-based rendering to terminal

### Styling

- **Colors**: Hex colors, 256-color palette, true color support
- **Borders**: Rounded, double, simple, and custom borders
- **Text Templates**: Inline styling with `[color]text[]`, `[bold]text[]`
- **Cell Styles**: Foreground, background, bold, italic, underline

## Examples

Run examples from the `examples/` directory:

```bash
cd examples
./run.sh simple       # Basic static demo
./run.sh counter      # Timer-based updates
./run.sh program-demo # Interactive text editor
./run.sh demo         # Original manual demo
./run.sh emoji-test   # Emoji rendering test
```

## Architecture

### Before (Manual)
```go
// Setup terminal manually
restore, _ := terminal.RawMode()
defer restore()
ansi.Print(ansi.EnterAltScreen)
defer ansi.Print(ansi.ExitAltScreen)

// Manual rendering
r := velvet.NewRenderer(os.Stdout, w, h)
b := buffer.NewGrid(w, h)
// ... build buffer
r.Write(b)
r.Flush()
```

### After (Program API)
```go
// Everything handled automatically!
p := velvet.NewProgram(model{})
p.Run()
```

## Template Syntax

```go
"[bold]Bold text[]"
"[red]Red text[]"
"[bold,cyan]Bold cyan text[]"
"[#ff5500]Custom hex color[]"
"Normal [green]green[] normal"
```

Supported styles:
- Colors: `red`, `green`, `blue`, `cyan`, `magenta`, `yellow`, `white`, `black`
- Modifiers: `bold`, `italic`, `underline`, `dim`
- Hex colors: `#rrggbb`

## Key Features

### Event-Driven Updates
Messages flow through your `Update` function, making state management predictable and testable.

### Buffer-Based Rendering
Unlike string-based TUIs, Velvet works with cell grids for pixel-perfect control over every character's appearance.

### Automatic Diffing
Only changed cells are redrawn, making updates efficient even for large terminals.

### Window Resize
Automatically handles `SIGWINCH` signals and sends `WindowSizeMsg` to your model.

### Command System
Run async operations that return messages to update your model:

```go
func loadData() velvet.Cmd {
    return func() velvet.Msg {
        data := fetchAPI()
        return DataLoadedMsg{data}
    }
}
```

## Comparison with Other Libraries

| Feature | Velvet | Bubbletea | Termbox |
|---------|--------|-----------|---------|
| Architecture | Buffer-based | String-based | Cell-based |
| Event Loop | Built-in | Built-in | Manual |
| Styling | Rich templates | External libs | Basic |
| Layout | Grid + alignment | Manual | Manual |
| Components | Built-in | External libs | None |

## Contributing

Contributions welcome! This is an active project.

## License

MIT

## Credits

Inspired by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and built with modern Go practices.

