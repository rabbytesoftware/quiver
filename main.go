package main

import (
	server "rounds.com.ar/watcher/packages/server"
	ui "rounds.com.ar/watcher/view"
	logger "rounds.com.ar/watcher/view/logger"
	api "rounds.com.ar/watcher/rest/api"
	packages_global_variables "rounds.com.ar/watcher/rest/shared/utils/packages/global-variables"
)

func main() {
	logger.It = logger.NewLogger()

	ui.Welcome(logger.It)

	logger.It.Load("api-main", "Loading packages...")
	packagesDir := "./pkgs"
	pkgServer := server.NewPackagesServer(packagesDir)

	if err := pkgServer.Discover(); err != nil {
		logger.It.Fatal("api-main", "Failed to discover packages: %v", err)
		return
	}

	if len(pkgServer.Packages) == 0 {
		logger.It.Warn("api-main", "No packages found in %s", packagesDir)
	}

	packageNames := make([][]string, 0, len(pkgServer.Packages))
	for _, pkg := range pkgServer.Packages {
		packageNames = append(packageNames, []string{pkg.Name, pkg.Version, pkg.URL, pkg.BuildNumber})
	}

	ui.Table("Packages", []string{"Name", "Version", "URL", "Build Number"}, packageNames)

	logger.It.Ok("api-main", "Packages loaded successfully")

	// Server API
	serverApi := api.CreateServerAPI(":8080")

	// Assign to packages global variable
	packages_global_variables.Packages = pkgServer.Packages 

	// Return error if server fails
	if err := serverApi.Run(); err != nil {
		logger.It.Error("api-main", "Error running server.")
	}
}
