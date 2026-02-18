# Velvet Program API

The Velvet Program API provides a bubbletea-style architecture for building terminal user interfaces with buffer-based rendering.

## Overview

The Program struct manages the application lifecycle, including:
- Terminal setup (raw mode, alternate screen)
- Input handling (keyboard events)
- Message/command processing
- Automatic rendering
- Window resize handling

## Core Concepts

### Model

Your application implements the `Model` interface:

```go
type Model interface {
    Init() Cmd                    // Called when program starts
    Update(Msg) (Model, Cmd)      // Called on each message
    View() *buffer.Grid           // Returns the buffer to render
}
```

### Messages

Messages are events sent to your Model's Update function:

- `KeyMsg` - Keyboard input
- `WindowSizeMsg` - Terminal resize
- `TickMsg` - Timer events
- `QuitMsg` - Program quit signal
- Custom messages from your commands

### Commands

Commands are functions that return messages:

```go
type Cmd func() Msg
```

Built-in commands:
- `Tick(duration)` - Send a TickMsg after a duration
- `Batch(cmds...)` - Execute multiple commands concurrently
- `Sequence(cmds...)` - Execute commands in order
- `Quit()` - Quit the program

## Basic Example

```go
package main

import (
    "github.com/ProductionPanic/velvet"
    "github.com/ProductionPanic/velvet/buffer"
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
        m.width = msg.Width
        m.height = msg.Height
    
    case velvet.KeyMsg:
        if msg.Type == velvet.KeyCtrlC {
            return m, func() velvet.Msg {
                return velvet.Quit()
            }
        }
    }
    
    return m, nil
}

func (m model) View() *buffer.Grid {
    b := buffer.NewGrid(m.width, m.height)
    // ... render your UI to the buffer
    return b
}

func main() {
    p := velvet.NewProgram(model{})
    if err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Key Types

### KeyMsg

```go
type KeyMsg struct {
    Type  KeyType
    Runes []rune
    Alt   bool
}
```

Key types include:
- `KeyRunes` - Regular characters
- `KeyEnter`, `KeyBackspace`, `KeyDelete`
- `KeyUp`, `KeyDown`, `KeyLeft`, `KeyRight`
- `KeyTab`, `KeyShiftTab`
- `KeyHome`, `KeyEnd`, `KeyPgUp`, `KeyPgDown`
- `KeyEsc`, `KeyCtrlC`, `KeyCtrlD`, `KeyCtrlZ`
- `KeySpace`

## Program Options

Configure the program with functional options:

```go
p := velvet.NewProgram(
    model{},
    velvet.WithAltScreen(true),      // Enable alternate screen
    velvet.WithInput(customReader),   // Custom input source
    velvet.WithOutput(customWriter),  // Custom output destination
    velvet.WithoutInput(),            // Disable input (for display-only)
)
```

## Examples

Check the `examples/` directory:
- `program-demo/` - Interactive text input example
- `counter/` - Timer-based counter with keyboard input

## Differences from Bubbletea

1. **Buffer-based rendering**: Instead of returning strings, `View()` returns `*buffer.Grid`
2. **Cell-level control**: Work with individual cells for precise styling and layout
3. **Built-in components**: Use `Container` for borders, padding, and text wrapping
4. **Native grid system**: Place components with alignment controls

## Advanced Usage

### Custom Commands

Create commands that return custom messages:

```go
type DataLoadedMsg struct {
    Data string
}

func loadData() velvet.Cmd {
    return func() velvet.Msg {
        // Do async work
        data := fetchFromAPI()
        return DataLoadedMsg{Data: data}
    }
}

func (m model) Init() velvet.Cmd {
    return loadData()
}
```

### Batch Commands

Execute multiple commands at once:

```go
func (m model) Init() velvet.Cmd {
    return velvet.Batch(
        loadData(),
        velvet.Tick(time.Second),
    )
}
```

### Timer Commands

Create repeating timers:

```go
func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
    switch msg := msg.(type) {
    case velvet.TickMsg:
        m.counter++
        return m, velvet.Tick(time.Second) // Keep ticking
    }
    return m, nil
}
```

