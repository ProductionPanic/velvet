# Quick Start Guide

## Running Examples

```bash
# From examples directory
./run.sh simple       # Basic static demo
./run.sh counter      # Timer-based counter
./run.sh program-demo # Interactive text input
```

## Creating a New Program

### 1. Define your model
```go
type model struct {
    width, height int
    // your state here
}
```

### 2. Implement Init()
```go
func (m model) Init() velvet.Cmd {
    return nil // or return a command
}
```

### 3. Implement Update()
```go
func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
    switch msg := msg.(type) {
    case velvet.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case velvet.KeyMsg:
        if msg.Type == velvet.KeyCtrlC {
            return m, func() velvet.Msg { return velvet.Quit() }
        }
    }
    return m, nil
}
```

### 4. Implement View()
```go
func (m model) View() *buffer.Grid {
    b := buffer.NewGrid(m.width, m.height)
    // render your UI
    return b
}
```

### 5. Run it
```go
func main() {
    p := velvet.NewProgram(model{})
    if err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Common Patterns

### Quit on Escape or Ctrl+C
```go
case velvet.KeyMsg:
    if msg.Type == velvet.KeyCtrlC || msg.Type == velvet.KeyEsc {
        return m, func() velvet.Msg { return velvet.Quit() }
    }
```

### Timer Updates
```go
func (m model) Init() velvet.Cmd {
    return velvet.Tick(time.Second)
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
    switch msg.(type) {
    case velvet.TickMsg:
        m.counter++
        return m, velvet.Tick(time.Second) // keep ticking
    }
    return m, nil
}
```

### Centered Container
```go
func (m model) View() *buffer.Grid {
    b := buffer.NewGrid(m.width, m.height)
    
    container := velvet.NewContainer().
        Padding(2).
        Border(style.RoundedBorder()).
        Width(m.width / 2)
    
    b.Place(container.Render("content"), buffer.AlignCenter, buffer.AlignMiddle)
    return b
}
```

### Handle Text Input
```go
case velvet.KeyMsg:
    switch msg.Type {
    case velvet.KeyRunes:
        m.input += string(msg.Runes)
    case velvet.KeyBackspace:
        if len(m.input) > 0 {
            m.input = m.input[:len(m.input)-1]
        }
    case velvet.KeyEnter:
        // process input
    }
```

## Key Types Reference

- `KeyRunes` - Regular text (check `msg.Runes`)
- `KeyEnter` - Enter/Return
- `KeyBackspace` - Backspace
- `KeyDelete` - Delete
- `KeyUp/Down/Left/Right` - Arrow keys
- `KeyTab`, `KeyShiftTab` - Tab keys
- `KeyHome`, `KeyEnd` - Home/End
- `KeyPgUp`, `KeyPgDown` - Page up/down
- `KeyEsc` - Escape
- `KeyCtrlC`, `KeyCtrlD`, `KeyCtrlZ` - Ctrl combinations
- `KeySpace` - Spacebar

## Program Options

```go
velvet.NewProgram(
    model{},
    velvet.WithAltScreen(false),     // Disable alt screen
    velvet.WithOutput(customWriter),  // Custom output
    velvet.WithInput(customReader),   // Custom input
    velvet.WithoutInput(),            // No input (display only)
)
```

