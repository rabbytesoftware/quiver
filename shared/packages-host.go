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

	pb "rounds.com.ar/sdk/package"
)

// PackagesHost manages the lifecycle of packages
type PackagesHost struct {
	PackagesDir string
	Packages    map[string]*Package
	NextPort    int
	mutex       sync.Mutex
}

// NewPackagesHost creates a new package host
func NewPackagesHost(packagesDir string) *PackagesHost {
	return &PackagesHost{
		PackagesDir: packagesDir,
		Packages:    make(map[string]*Package),
		NextPort:    50051, // Starting port for packages
	}
}

// DiscoverPackages finds all executable packages in the packages directory
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
		if isExecutable(filePath) {
			// Generate a port for this package
			port := h.NextPort
			h.NextPort++

			packagePath := filePath
			h.Packages[packagePath] = &Package{
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
	var wg sync.WaitGroup
	var list []struct {
		Item     string
		Callback func() (bool, error)
	}

	for path, info := range h.Packages {
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

	wg.Wait()
	return list, nil
}

// CloseAllPackages stops all packages and cleans up
func (h *PackagesHost) CloseAllPackages() {
	for _, info := range h.Packages {
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