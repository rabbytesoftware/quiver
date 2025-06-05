package main

import (
	// server "github.com/rabbytesoftware/quiver/packages/server"

	logger "github.com/rabbytesoftware/quiver/logger"
	packages "github.com/rabbytesoftware/quiver/packages"
	ui "github.com/rabbytesoftware/quiver/view"
)

func main() {
	logs := logger.NewLogger("ini")

	ui.Welcome()

  	// logs.Load("Loading Netbridge...")
	// netbridge, err := netbridge.NewNetbridge()
	// if err != nil {
	// 	logs.Fatal("Failed to init Netbridge: %v", err)
	// 	return
	// }
	// logs.Ok("Quiver instance registered with IP: %s", netbridge.PublicIP)

	logs.Load("Loading packages...")
	packagesDir := "./pkgs"
	arrowsServer := packages.NewArrowsServer(packagesDir)

	arrowsServer.Load("./template/cs2.yaml")

	// if err := pkgServer.Discover(); err != nil {
	// 	logs.Warn("Failed to discover packages: %v", err)
	// }

	// if len(pkgServer.Packages) == 0 {
	// 	logs.Warn("No packages found in %s", packagesDir)
	// }

	// packageNames := make([][]string, 0, len(pkgServer.Packages))
	// for _, pkg := range pkgServer.Packages {
	// 	packageNames = append(packageNames, []string{pkg.Name, pkg.Version, pkg.URL, pkg.BuildNumber})
	// }
	// ui.Table("Packages", []string{"Name", "Version", "URL", "Build Number"}, packageNames)
	// logs.Ok("Packages loaded successfully")

	// logs.Load("Loading API server...")
	// serverApi := api.CreateServerAPI(":8080", &pkgServer.Packages)
	// if err := serverApi.Run(); err != nil {
	// 	logs.Error("Error running API server")
	// }
}
