package api_main

import (
	"fmt"

	"rounds.com.ar/watcher/rest/api"
)

func main() {
	server := api.CreateServerAPI(":8080")

	if err := server.Run(); err != nil {
		fmt.Errorf("Error running server.")
	}
}