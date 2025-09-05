package main

import (
	"log"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui"
)

func main() {
	ui := ui.NewTUI()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
