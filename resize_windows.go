//go:build windows

package velvet

import (
	"time"

	"github.com/ProductionPanic/velvet/terminal"
)

func handleResize(p *Program) {
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				w, h, err := terminal.GetSize()
				if err == nil && (w != p.width || h != p.height) {
					p.Send(WindowSizeMsg{Width: w, Height: h})
				}
			}
		}
	}()
}
