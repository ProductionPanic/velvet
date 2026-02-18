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
        if msg.String() == "ctrl+c" {
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
    switch msg.String() {
    case "ctrl+c", "esc":
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
    switch msg.String() {
    case "backspace":
        if len(m.input) > 0 {
            m.input = m.input[:len(m.input)-1]
        }
    case "enter":
        // process input
    default:
        // Regular character input
        if len(msg.Runes) > 0 {
            m.input += string(msg.Runes)
        }
    }
```

## Key Reference

The `KeyMsg.String()` method returns string representations of keys that can be used in switch statements:

**Special Keys:**
- `"enter"` - Enter/Return key
- `"backspace"` - Backspace key
- `"delete"` - Delete key
- `"esc"` - Escape key
- `"tab"` - Tab key
- `"shift+tab"` - Shift+Tab key
- `" "` - Spacebar (literal space character)

**Arrow Keys:**
- `"up"`, `"down"`, `"left"`, `"right"` - Arrow keys

**Navigation Keys:**
- `"home"`, `"end"` - Home/End keys
- `"pgup"`, `"pgdown"` - Page up/down keys

**Control Keys:**
- `"ctrl+c"`, `"ctrl+d"`, `"ctrl+z"` - Ctrl combinations

**Regular Characters:**
- Single characters like `"a"`, `"5"`, `"@"` - Regular text input
- UTF-8 characters and emoji are supported

**Alt Combinations:**
- Alt key combinations return `"alt+X"` format (e.g., `"alt+a"`)

**Accessing Raw Input:**
- Use `msg.Runes` to get the raw rune slice for text input
- Use `msg.Alt` to check if Alt modifier is pressed
- Use `msg.Matches("key1", "key2", ...)` helper to check multiple keys

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

