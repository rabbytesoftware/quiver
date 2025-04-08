package shared

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pb "rounds.com.ar/watcher/sdk/package"
	shared "rounds.com.ar/watcher/shared"
)

// PackagesHost manages the lifecycle of packages
type PackagesHost struct {
	PackagesDir string
	Packages    map[string]*shared.Package
	NextPort    int
	mutex       sync.Mutex
}

// NewPackagesHost creates a new package host
func NewPackagesHost(packagesDir string) *PackagesHost {
	return &PackagesHost{
		PackagesDir:     packagesDir,
		Packages:        make(map[string]*shared.Package),
		NextPort:        50051, // TODO: With the NetBridge system, we can check if 
								// TODO: the port is on use by another software, and
								// TODO: try another.
	}
}

// DiscoverPackages finds all executable packages and watcher packages in the packages directory
func (h *PackagesHost) DiscoverPackages() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	files, err := ioutil.ReadDir(h.PackagesDir)
	if err != nil {
		return fmt.Errorf("failed to read packages directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(h.PackagesDir, file.Name())

		// Check if it's a watcher package (.watcher extension)
		if strings.ToLower(filepath.Ext(filePath)) == ".watcher" {
			// Extract the watcher package
			watcherPkg, err := shared.ExtractWatcherPackage(filePath)
			if err != nil {
				fmt.Printf("Warning: Failed to extract watcher package %s: %v\n", filePath, err)
				continue
			}
			
			// Generate a port for this package
			port := h.NextPort
			h.NextPort++
			
			// Store the watcher package
			watcherPkg.BasePort = port
			h.Packages[filePath] = watcherPkg

			// Cleanup the watcher package
			shared.CleanupWatcherPackage(watcherPkg)
			
			continue
		}
	}

	return nil
}

// CloseAllPackages stops all packages and cleans up
func (h *PackagesHost) CloseAllPackages() {
	for path, pkg := range h.Packages {
		if pkg.Connection != nil {
			// Try to stop the package gracefully
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			pkg.Client.Exit(ctx, &pb.Empty{})
			cancel()

			// Close the connection
			pkg.Connection.Close()
		}

		// Kill the process if it's still running
		if pkg.Process != nil {
			pkg.Process.Kill()
		}
		
		// Clean up watcher package if applicable
		if wp, ok := h.Packages[path]; ok {
			shared.CleanupWatcherPackage(wp)
		}
	}
}
