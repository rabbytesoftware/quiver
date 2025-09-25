package main

import (
	"fmt"

	"quiver.core/internal/core/metadata"
)

func main() {
	fmt.Println("UUID:", metadata.GetUUID())
	fmt.Println("Version:", metadata.GetVersion())
}
