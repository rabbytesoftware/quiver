package main

import (
	"log"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/tui"
)

func main() {
	logger, err := tui.RunTUI()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Hello world")
}
