package velvet

type ProgramOption func(*Program)

type programOptions struct {
	AltScreen bool
}

func WithAltScreen() ProgramOption {
	return func(p *Program) {
		p.programOptions.AltScreen = true
	}
}

func SetWindowTitle(title string) ProgramOption {
	return func(p *Program) {
		p.programTitle = title
	}
}
