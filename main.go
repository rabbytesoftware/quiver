package main

import (
	netbridge "github.com/rabbytesoftware/quiver/netbridge"
	server "github.com/rabbytesoftware/quiver/packages/server"

	ui "github.com/rabbytesoftware/quiver/view"
	logger "github.com/rabbytesoftware/quiver/view/logger"
)

func main() {
	logger.It = logger.NewLogger()

	ui.Welcome(logger.It)

	logger.It.Load("Loading Netbridge...")
	netbridge, err := netbridge.NewNetbridge()
	if err != nil {
		logger.It.Fatal("Failed to init Netbridge: %v", err)
		return
	}
	logger.It.Ok("Device registered with IP: %s", netbridge.PublicIP)

	logger.It.Load("Loading packages...")
	packagesDir := "./pkgs"
	pkgServer := server.NewPackagesServer(packagesDir)

	if err := pkgServer.Discover(); err != nil {
		logger.It.Fatal("Failed to discover packages: %v", err)
		return
	}

	if len(pkgServer.Packages) == 0 {
		logger.It.Warn("No packages found in %s", packagesDir)
	}

	packageNames := make([][]string, 0, len(pkgServer.Packages))
	for _, pkg := range pkgServer.Packages {
		packageNames = append(packageNames, []string{pkg.Name, pkg.Version, pkg.URL, pkg.BuildNumber})
	}

	ui.Table("Packages", []string{"Name", "Version", "URL", "Build Number"}, packageNames)

	logger.It.Ok("Packages loaded successfully")

	// Server API
	serverApi := api.CreateServerAPI(":8080")

	// Assign to packages global variable
	packages_global_variables.Packages = pkgServer.Packages 

	// Return error if server fails
	if err := serverApi.Run(); err != nil {
		logger.It.Error("Error running server.")
	}
}
