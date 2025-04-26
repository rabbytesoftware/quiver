package main

import (
	netbridge "rounds.com.ar/watcher/netbridge"
	server "rounds.com.ar/watcher/packages/server"

	ui "rounds.com.ar/watcher/view"
	logger "rounds.com.ar/watcher/view/logger"
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
}
