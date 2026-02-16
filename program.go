package velvet

import (
	"context"
	"io"
	"os"

	"golang.org/x/term"
)

type Model interface {
	Init() Msg
	Update(Msg) (Model, Cmd)
	View() *Buffer
}

type Program struct {
	initialModel Model

	programOptions programOptions // struct of all options for the program, set by the user when creating the program
	programTitle   string         // title of the terminal

	externalCtx context.Context // context that can be passed in by the user

	ctx    context.Context    // context for the program
	cancel context.CancelFunc // function to cancel the program's context

	// channels for passing data
	msgs     chan Msg
	errs     chan error
	finished chan struct{}

	output io.Writer
	input  io.Reader

	renderer renderer
}

func NewProgram(model Model, options ...ProgramOption) *Program {
	p := &Program{
		initialModel: model,
		msgs:         make(chan Msg),
	}

	for _, opt := range options {
		opt(p)
	}

	if p.externalCtx == nil {
		p.externalCtx = context.Background()
	}

	p.ctx, p.cancel = context.WithCancel(p.externalCtx)

	if p.output == nil {
		p.output = os.Stdout
	}

	if p.input == nil {
		p.input = os.Stdin
	}

	return p
}

func (p *Program) Run() (Model, error) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}
	if p.renderer == (renderer{}) {
		p.renderer = newRenderer(p.output, w, h)
	}

}
