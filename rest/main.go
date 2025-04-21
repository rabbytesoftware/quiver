package main

import (
	server "rounds.com.ar/watcher/packages/server"
	api "rounds.com.ar/watcher/rest/api"
	packages_global_variables "rounds.com.ar/watcher/rest/shared/utils/packages/global-variables"
	logger "rounds.com.ar/watcher/view/logger"
)


func main() {
	serverApi := api.CreateServerAPI(":8080")

	logger.It.Load("Loading packages...")
	packagesDir := "../pkgs"
	pkgServer := server.NewPackagesServer(packagesDir)

	// Get packages list
	if err := pkgServer.Discover(); err != nil {
		logger.It.Fatal("failed to discover packages: %v", err)
		return
	}
	
	// Check if any result
	if len(pkgServer.Packages) == 0 {
		logger.It.Warn("no packages found in %s", packagesDir)
		return
	}

	// Assign to global variable
	packages_global_variables.Packages = pkgServer.Packages 

	if err := serverApi.Run(); err != nil {
		logger.It.Error("Error running server.")
	}
}