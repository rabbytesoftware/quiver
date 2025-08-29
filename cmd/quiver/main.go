package main

import (
	"log"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui"
)

func main() {
	ui := ui.NewUI()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
