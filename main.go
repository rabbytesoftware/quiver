package main

import (
	"log"

	host "rounds.com.ar/watcher/shared/host"
	view "rounds.com.ar/watcher/view"
)

func main() {
	view.Init()

	// Create package host
	packagesDir := "./pkgs"
	host := host.NewPackagesHost(packagesDir)

	// Discover packages
	if err := host.DiscoverPackages(); err != nil {
		log.Fatalf("Failed to discover packages: %v", err)
	}

	if len(host.Packages) == 0 {
		log.Printf("No packages found in %s", packagesDir)
		return
	}

	// Display loaded packages
	packageNames := make([][]string, 0, len(host.Packages))
	for _, pkg := range host.Packages {
		meta := pkg.Metadata
		packageNames = append(packageNames, []string{meta.Name, meta.Version, meta.URL, meta.BuildNumber})
	}

	view.Table("Packages", []string{"Name", "Version", "URL", "Build Number"}, packageNames)
}
