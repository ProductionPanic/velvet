package terminal

import (
	"os"

	"golang.org/x/term"
)

func GetSize() (width, height int, err error) {
	return term.GetSize(
		int(os.Stdout.Fd()),
	)
}

func RawMode() (restore func(), err error) {
	oldState, err := term.MakeRaw(
		int(os.Stdin.Fd()),
	)
	if err != nil {
		return nil, err
	}

	restore = func() {
		_ = term.Restore(
			int(os.Stdin.Fd()),
			oldState,
		)
	}

	return restore, nil
}
