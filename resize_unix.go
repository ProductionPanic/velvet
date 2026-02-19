//go:build !windows

package velvet

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ProductionPanic/velvet/terminal"
)

func handleResize(p *Program) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	go func() {
		for {
			select {
			case <-p.quit:
				return
			case <-sigChan:
				w, h, err := terminal.GetSize()
				if err == nil {
					p.Send(WindowSizeMsg{Width: w, Height: h})
				}
			}
		}
	}()
}
