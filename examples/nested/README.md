# Nested Components Example

This example demonstrates how the `Drawable` interface enables flexible, composable component hierarchies in velvet.

## What This Example Shows

### Before (Without Drawable Interface)
Previously, `View()` returned `*buffer.Grid` directly, which meant:
- Layout containers could only work with `*Component` types
- Building complex nested layouts required manual grid manipulation
- Custom drawable types couldn't participate in layouts
- Type safety was limited

### After (With Drawable Interface)
Now, `Render()` returns `Drawable`, which means:
- **Any type can be renderable** - As long as it implements `Render() *buffer.Grid`
- **Flexible composition** - Layout containers accept `[]Drawable`, enabling arbitrary nesting
- **Reusable component factories** - Functions can return `Drawable` for easy composition
- **Type-safe hierarchies** - The compiler ensures all rendered types implement the interface

## Example Structure

```
model.Render() → Drawable (FlexContainer)
    ├── statsSection → Drawable (FlexContainer)
    │   ├── cpuCard → Drawable (Component)
    │   ├── memCard → Drawable (Component)
    │   └── diskCard → Drawable (Component)
    └── bottomSection → Drawable (FlexContainer)
        └── networkCard → Drawable (Component)
```

## Key Benefits Demonstrated

### 1. Composable Component Factories
```go
func (m model) createStatCard(label string, value int, accentHex string) velvet.Drawable {
    return velvet.NewComponent().
        // ... configuration ...
        Content(content)
}
```

The method returns `Drawable`, making it easy to compose with layouts:
```go
cpuCard := m.createStatCard("CPU", m.stats.CPUUsage, "#ff6b6b")
memCard := m.createStatCard("Memory", m.stats.MemoryUsage, "#4ecdc4")

statsSection := velvet.NewFlexContainer().
    SetSize(m.width, m.height/2).
    Add(cpuCard, memCard)  // ✓ Works because both are Drawable
```

### 2. Arbitrary Nesting Depth
Layouts can contain other layouts, which contain components, which could contain more layouts:
```go
mainLayout := velvet.NewFlexContainer().
    SetSize(m.width, m.height).
    Vertical().
    Add(statsSection, bottomSection)  // Both are FlexContainers (Drawable)
```

### 3. Type Safety
The compiler ensures composition is correct:
```go
// This works ✓
container.Add(component, layout, customDrawable)

// This would fail ✗ (caught at compile time)
container.Add("not a drawable")
```

## Running the Example

```bash
cd examples/nested
go build -o nested
./nested
```

Press `ESC` or `Ctrl+C` to exit.

## Extension Ideas

With the `Drawable` interface, you could now easily:

1. **Create custom Drawable types** - Implement `Render() *buffer.Grid` for specialized components
2. **Build container combinations** - Mix FlexContainer and GridContainer at any depth
3. **Prepare for interactivity** - The interface is positioned to support mouse event routing to nested components (future enhancement)
4. **Reuse layout logic** - Drawable factories can be shared across projects

## Why This Matters for Future Development

The `Drawable` interface sets up the architecture for:
- **Position tracking** - Each drawable will eventually track its bounding box
- **Mouse event routing** - Click events can be delegated from parent layouts to nested children
- **Accessibility** - Nested components can report their positions and interactions
- **Advanced compositions** - Custom layout systems can be built on top of the interface

This refactor transforms velvet from a component-focused library to a fully composable UI system.

