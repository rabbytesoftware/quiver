package shared

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
	// Store extracted watcher packages
	WatcherPackages map[string]*shared.WatcherPackage
}

// NewPackagesHost creates a new package host
func NewPackagesHost(packagesDir string) *PackagesHost {
	return &PackagesHost{
		PackagesDir:     packagesDir,
		Packages:        make(map[string]*shared.Package),
		WatcherPackages: make(map[string]*shared.WatcherPackage),
		NextPort:        50051, // Starting port for packages
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
			h.WatcherPackages[filePath] = watcherPkg
			
			// Register the package in the packages map
			h.Packages[filePath] = &shared.Package{
				BasePort: port,
				
			}
			
			continue
		}
		
		// Handle regular executable packages (non-watcher)
		if isExecutable(filePath) {
			// Generate a port for this package
			port := h.NextPort
			h.NextPort++

			packagePath := filePath
			h.Packages[packagePath] = &shared.Package{
				BasePort: port,
			}
		}
	}

	return nil
}

func (h *PackagesHost) GetAllPackages() (
	[]struct {
		Item     string
		Callback func() (bool, error)
	},
	error,
) {
	var list []struct {
		Item     string
		Callback func() (bool, error)
	}

	for path, info := range h.Packages {
		path := path // Create local copy for closure
		info := info // Create local copy for closure
		
		list = append(list, struct {
			Item     string
			Callback func() (bool, error)
		}{
			Item: path,
			Callback: func() (bool, error) {
				err := h.startPackage(path, info)
				return err != nil, err
			},
		})
	}

	return list, nil
}

// CloseAllPackages stops all packages and cleans up
func (h *PackagesHost) CloseAllPackages() {
	for path, info := range h.Packages {
		if info.Connection != nil {
			// Try to stop the package gracefully
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			info.Client.Exit(ctx, &pb.Empty{})
			cancel()

			// Close the connection
			info.Connection.Close()
		}

		// Kill the process if it's still running
		if info.Process != nil {
			info.Process.Kill()
		}
		
		// Clean up watcher package if applicable
		if wp, ok := h.WatcherPackages[path]; ok {
			shared.CleanupWatcherPackage(wp)
		}
	}
}

// isExecutable checks if a file is executable based on the OS
func isExecutable(path string) bool {
	if runtime.GOOS == "windows" {
		// On Windows, check if the file has .exe extension
		return strings.ToLower(filepath.Ext(path)) == ".exe"
	} else {
		// On Unix-like systems, check file permissions
		info, err := os.Stat(path)
		if err != nil {
			return false
		}
		return info.Mode()&0111 != 0
	}
}