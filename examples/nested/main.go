package main

import (
	"fmt"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

// This example demonstrates how the Drawable interface enables flexible nesting
// of components and layout containers. We build a dashboard with multiple sections,
// each containing their own nested components.

type model struct {
	width, height int
	stats         Stats
}

type Stats struct {
	CPUUsage    int
	MemoryUsage int
	DiskUsage   int
	NetworkIn   int64
	NetworkOut  int64
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
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, velvet.Quit
		}
	}

	return m, nil
}

// Render demonstrates nested Drawable composition
// We build a hierarchy: Layout > Sections > Cards > Content
func (m model) Render() velvet.Drawable {
	// Create individual stat cards - these are Components (Drawables)
	cpuCard := m.createStatCard("CPU", m.stats.CPUUsage, "#ff6b6b")
	memCard := m.createStatCard("Memory", m.stats.MemoryUsage, "#4ecdc4")
	diskCard := m.createStatCard("Disk", m.stats.DiskUsage, "#45b7d1")

	// Create a FlexContainer (also a Drawable) that lays out the cards horizontally
	statsSection := velvet.NewFlexContainer().
		SetSize(m.width, m.height/2).
		Horizontal().
		Add(cpuCard, memCard, diskCard)

	// Create network info card
	networkCard := m.createNetworkCard()

	// Create another FlexContainer for the bottom section
	// This demonstrates that we can nest Drawables arbitrarily deep
	bottomSection := velvet.NewFlexContainer().
		SetSize(m.width, m.height/2).
		Vertical().
		Add(networkCard)

	// Create the main layout that combines both sections vertically
	mainLayout := velvet.NewFlexContainer().
		SetSize(m.width, m.height).
		Vertical().
		Add(statsSection, bottomSection)

	return mainLayout
}

// createStatCard creates a reusable stat card component
// Notice how this returns a Drawable (Component), making it composable
func (m model) createStatCard(label string, value int, accentHex string) velvet.Drawable {
	content := fmt.Sprintf("[%s]%s[]\n\n[bold]%d%%[]", accentHex, label, value)

	return velvet.NewComponent().
		Foreground(color.FromHex("#ffffff")).
		Background(color.FromHex("#1a1a1a")).
		BorderColor(color.FromHex(accentHex)).
		Border(style.RoundedBorder()).
		Padding(1, 2).
		Content(content)
}

// createNetworkCard creates a more complex card with multiple lines
func (m model) createNetworkCard() velvet.Drawable {
	content := fmt.Sprintf(
		"[bold cyan]Network Activity[]\n\n"+
			"[dim]Upload:[]\n  %d KB/s\n\n"+
			"[dim]Download:[]\n  %d KB/s",
		m.stats.NetworkOut/1024,
		m.stats.NetworkIn/1024,
	)

	return velvet.NewComponent().
		Foreground(color.FromHex("#ffffff")).
		Background(color.FromHex("#1a1a1a")).
		BorderColor(color.FromHex("#a29bfe")).
		Border(style.RoundedBorder()).
		Padding(1, 2).
		Content(content)
}

func initialModel() model {
	return model{
		stats: Stats{
			CPUUsage:    65,
			MemoryUsage: 42,
			DiskUsage:   78,
			NetworkIn:   5120000,
			NetworkOut:  1024000,
		},
	}
}

func main() {
	p := velvet.NewProgram(initialModel())
	if err := p.Run(); err != nil {
		panic(err)
	}
}
