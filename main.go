package main

import (
	"log"

	server "rounds.com.ar/watcher/packages/server"
	view "rounds.com.ar/watcher/view"
)

func main() {
	view.Init()

	// Create package host
	packagesDir := "./pkgs"
	pkgServer := server.NewPackagesServer(packagesDir)

	// Discover packages
	if err := pkgServer.Discover(); err != nil {
		log.Fatalf("Failed to discover packages: %v", err)
	}

	if len(pkgServer.Packages) == 0 {
		log.Printf("No packages found in %s", packagesDir)
		return
	}

	// Display loaded packages
	packageNames := make([][]string, 0, len(pkgServer.Packages))
	for _, pkg := range pkgServer.Packages {
		packageNames = append(packageNames, []string{pkg.Name, pkg.Version, pkg.URL, pkg.BuildNumber})
	}

	view.Table("Packages", []string{"Name", "Version", "URL", "Build Number"}, packageNames)
}
