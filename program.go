package velvet

import (
	"io"
	"os"

	"github.com/ProductionPanic/velvet/ansi"
	"github.com/ProductionPanic/velvet/terminal"
)

// Model is the interface that user programs implement
type Model interface {
	// Init is called when the program starts
	Init() Cmd

	// Update is called when a message is received
	Update(Msg) (Model, Cmd)

	// Render returns a Drawable to render
	Render() Drawable
}

// Program manages the application lifecycle
type Program struct {
	model         Model
	renderer      *Renderer
	output        io.Writer
	input         io.Reader
	msgs          chan Msg
	cmds          chan Cmd
	quit          chan struct{}
	width, height int

	// Options
	mouseAllMotion bool
	inputEnabled   bool
}

// ProgramOption is a functional option for configuring a Program
type ProgramOption func(*Program)

// WithMouseAllMotion enables mouse motion events
func WithMouseAllMotion(enable bool) ProgramOption {
	return func(p *Program) {
		p.mouseAllMotion = enable
	}
}

// WithInput sets a custom input reader
func WithInput(r io.Reader) ProgramOption {
	return func(p *Program) {
		p.input = r
	}
}

// WithOutput sets a custom output writer
func WithOutput(w io.Writer) ProgramOption {
	return func(p *Program) {
		p.output = w
	}
}

// WithoutInput disables input reading
func WithoutInput() ProgramOption {
	return func(p *Program) {
		p.inputEnabled = false
	}
}

// NewProgram creates a new program with the given model
func NewProgram(model Model, opts ...ProgramOption) *Program {
	p := &Program{
		model:        model,
		output:       os.Stdout,
		input:        os.Stdin,
		msgs:         make(chan Msg),
		cmds:         make(chan Cmd),
		quit:         make(chan struct{}),
		inputEnabled: true,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Run starts the program
func (p *Program) Run() error {
	// Get terminal size
	w, h, err := terminal.GetSize()
	if err != nil {
		return err
	}
	p.width = w
	p.height = h

	// Setup terminal
	if p.inputEnabled {
		restore, err := terminal.RawMode()
		if err != nil {
			return err
		}
		defer restore()
	}

	// Alternate screen is required because the renderer uses absolute cursor positioning.
	// Without it, the terminal scrollback would be corrupted and previous content destroyed.
	ansi.Print(ansi.EnterAltScreen)
	defer ansi.Print(ansi.ExitAltScreen)

	if p.mouseAllMotion {
		ansi.Print(ansi.EnableMouseAllMotion)
	} else {
		ansi.Print(ansi.EnableMouseClick)
	}
	ansi.Print(ansi.EnableMouseSGRMouse)

	defer ansi.Print(ansi.DisableMouseAllMotion, ansi.DisableMouseClick, ansi.DisableMouseSGRMouse)

	ansi.Print(
		ansi.ClearScreen,
		ansi.SetCursorPosition(1, 1),
		ansi.HideCursor,
	)
	defer ansi.Print(ansi.ShowCursor, ansi.Reset)

	// Create renderer
	p.renderer = NewRenderer(p.output, p.width, p.height)

	// Handle window resize signals
	handleResize(p)

	// Start input reader
	if p.inputEnabled {
		go readInput(p.msgs, p.input)
	}

	// Start command processor
	go p.processCommands()

	// Send initial window size to model
	var cmd Cmd
	p.model, cmd = p.model.Update(WindowSizeMsg{Width: p.width, Height: p.height})
	if cmd != nil {
		p.cmds <- cmd
	}

	// Initialize model
	if initCmd := p.model.Init(); initCmd != nil {
		p.cmds <- initCmd
	}

	// Initial render
	p.render()

	// Event loop
	return p.eventLoop()
}

// eventLoop processes messages until quit
func (p *Program) eventLoop() error {
	for {
		select {
		case <-p.quit:
			return nil
		case msg := <-p.msgs:
			// Check for quit message
			if _, ok := msg.(QuitMsg); ok {
				close(p.quit)
				return nil
			}

			// Handle window resize
			if wsm, ok := msg.(WindowSizeMsg); ok {
				p.width = wsm.Width
				p.height = wsm.Height
				p.renderer = NewRenderer(p.output, p.width, p.height)
			}

			// Update model
			var cmd Cmd
			p.model, cmd = p.model.Update(msg)
			if cmd != nil {
				p.cmds <- cmd
			}

			// Render
			p.render()
		}
	}
}

// render draws the current view
func (p *Program) render() {
	drawable := p.model.Render()
	if drawable != nil {
		grid := drawable.Render()
		p.renderer.Write(grid)
		p.renderer.Flush()
	}
}

// processCommands runs commands in the background
func (p *Program) processCommands() {
	for {
		select {
		case <-p.quit:
			return
		case cmd := <-p.cmds:
			if cmd == nil {
				continue
			}

			// Execute command in goroutine
			go func(c Cmd) {
				msg := c()
				if msg != nil {
					// Handle batch messages
					if batch, ok := msg.(batchMsg); ok {
						for _, cmd := range batch {
							if cmd != nil {
								p.cmds <- cmd
							}
						}
						return
					}

					// Handle sequence messages
					if seq, ok := msg.(sequenceMsg); ok {
						for _, cmd := range seq {
							if cmd != nil {
								msg := cmd()
								if msg != nil {
									p.msgs <- msg
								}
							}
						}
						return
					}

					p.msgs <- msg
				}
			}(cmd)
		}
	}
}

// Send sends a message to the program
func (p *Program) Send(msg Msg) {
	p.msgs <- msg
}

// Quit quits the program
func (p *Program) Quit() {
	p.msgs <- QuitMsg{}
}
