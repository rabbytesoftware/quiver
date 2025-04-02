package main

import (
	"log"

	host "rounds.com.ar/watcher/shared"
	view "rounds.com.ar/watcher/view"
	websocket "rounds.com.ar/watcher/websocket"
)

func main() {
	view.Init()

	// Create package host
	packagesDir := "./pkgs"
	host := host.NewPackageHost(packagesDir)

	// Discover packages
	if err := host.DiscoverPackages(); err != nil {
		log.Fatalf("Failed to discover packages: %v", err)
	}

	if len(host.Packages) == 0 {
		log.Printf("No packages found in %s", packagesDir)
		return
	}

	// Load all packages -> Test the package host
	packages, err := host.GetAllPackages()
	if err != nil {
		log.Fatalf("Failed to load packages: %v", err)
	}
	view.ProgressLoader("Loading Packages", "Loading", packages)
	defer host.CloseAllPackages()

	// Display loaded packages
	packageNames := make([][]string, 0, len(host.Packages))
	for _, info := range host.Packages {
		packageNames = append(packageNames, []string{info.Name, info.Version, info.Description, info.Icon})
	}

	view.Table("Packages", []string{"Name", "Version", "Description", "Icon"}, packageNames)

	websocket.Init()
}
