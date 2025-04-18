package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	pkgs "rounds.com.ar/watcher/packages"
	logger "rounds.com.ar/watcher/view/logger"
)

type PackagesServer struct {
	Packages    map[string]*pkgs.Package

	PackagesDir string
	NextPort    int
	mutex       sync.Mutex
}

func NewPackagesServer(packagesDir string) *PackagesServer {
	return &PackagesServer{
		PackagesDir:     packagesDir,
		Packages:        make(map[string]*pkgs.Package),
		NextPort:        50051, // TODO: With the NetBridge system, we can check if 
								// TODO: the port is on use by another software, and
								// TODO: try another.
	}
}

// DiscoverPackages finds all executable packages and watcher packages in the packages directory
func (h *PackagesServer) Discover() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	files, err := os.ReadDir(h.PackagesDir)
	if err != nil {
		return fmt.Errorf("failed to read packages directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(h.PackagesDir, file.Name())

		// ? Check if it's a watcher package (.watcher extension)
		if strings.ToLower(filepath.Ext(filePath)) == ".watcher" {
			watcherPkg := &pkgs.Package{
				Runtime: &pkgs.PackageRuntime{
					PackagePath: filePath,
				},
			}

			watcherPkg.Runtime.BasePort = h.NextPort
			h.NextPort++

			err := watcherPkg.Init()
			if err != nil {
				logger.It.Warn("Failed to extract watcher package %s: %v\n", filePath, err)
				continue
			}

			err = watcherPkg.Shutdown()
			if err != nil {
				logger.It.Warn("Warning: Failed to remove watcher package %s: %v\n", filePath, err)
				continue
			}
			
			h.Packages[filePath] = watcherPkg
			continue
		}
	}

	return nil
}

// CloseAllPackages stops all packages and cleans up
func (h *PackagesServer) CloseAll() {
	for _, pkg := range h.Packages {
		pkg.Shutdown()
	}
}
