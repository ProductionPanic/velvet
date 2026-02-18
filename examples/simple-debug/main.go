package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ProductionPanic/velvet"
	"github.com/ProductionPanic/velvet/buffer"
	"github.com/ProductionPanic/velvet/style"
	"github.com/ProductionPanic/velvet/style/color"
)

type model struct {
	width, height int
	logger        *log.Logger
}

func (m model) Init() velvet.Cmd {
	m.logger.Printf("Init called")
	return nil
}

func (m model) Update(msg velvet.Msg) (velvet.Model, velvet.Cmd) {
	switch msg := msg.(type) {
	case velvet.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.logger.Printf("WindowSizeMsg: %dx%d", m.width, m.height)

	case velvet.KeyMsg:
		m.logger.Printf("KeyMsg: %v", msg)
		// Quit on Ctrl+C or Escape
		if msg.Type == velvet.KeyCtrlC || msg.Type == velvet.KeyEsc {
			return m, func() velvet.Msg {
				return velvet.Quit()
			}
		}
	}

	return m, nil
}

func (m model) View() *buffer.Grid {
	m.logger.Printf("View called: %dx%d", m.width, m.height)
	b := buffer.NewGrid(m.width, m.height)

	if m.width == 0 || m.height == 0 {
		m.logger.Printf("WARNING: Zero dimensions!")
		return b
	}

	component := velvet.NewContainer().
		Foreground(color.FromHex("#ffffff")).
		Padding(1).
		Border(style.RoundedBorder()).
		BorderColor(color.FromHex("#00ff00")).
		Width(m.width / 4)

	content := fmt.Sprintf(
		"[bold,green]Simple Example[]\n\n"+
			"Terminal size: [cyan]%dx%d[]\n\n"+
			"Hello [red]velvet![] This is working.\n\n"+
			"[dim]Press ESC or Ctrl+C to quit[]",
		m.width, m.height,
	)

	rendered := component.Render(content)
	m.logger.Printf("Rendered container: %dx%d", rendered.Width, rendered.Height)

	b.Place(rendered, buffer.AlignCenter, buffer.AlignMiddle)

	return b
}

func main() {
	// Setup logging to file
	logFile, err := os.OpenFile("/tmp/velvet-debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Println("Starting program")

	m := model{logger: logger}
	p := velvet.NewProgram(m)

	if err := p.Run(); err != nil {
		logger.Printf("Error: %v", err)
		panic(err)
	}

	logger.Println("Program ended")
}
