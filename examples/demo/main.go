package main

import (
	"os"

	"github.com/ProductionPanic/velvet"
	"golang.org/x/term"
)

func main() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	renderer := velvet.newRenderer(
		os.Stdout,
		w,
		h,
	)

}
