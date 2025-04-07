package shared

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	entrypoint "rounds.com.ar/watcher/sdk/base/entrypoint"
)

// ExtractWatcherPackage extracts a .watcher file to a temporary directory
func ExtractWatcherPackage(packagePath string) (*Package, error) {
	// Create a unique temporary directory
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("watcher_%s_%d", 
		filepath.Base(packagePath), os.Getpid()))
	
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	// Open the package file
	file, err := os.Open(packagePath)
	if err != nil {
		os.RemoveAll(tempDir) // Clean up on error
		return nil, fmt.Errorf("failed to open package file: %w", err)
	}
	defer file.Close()
	
	// Create gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		os.RemoveAll(tempDir) // Clean up on error
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()
	
	// Create tar reader
	tr := tar.NewReader(gr)
	
	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.RemoveAll(tempDir) // Clean up on error
			return nil, fmt.Errorf("error reading archive: %w", err)
		}
		
		targetPath := filepath.Join(tempDir, header.Name)
		
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				os.RemoveAll(tempDir) // Clean up on error
				return nil, fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create file
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				os.RemoveAll(tempDir) // Clean up on error
				return nil, fmt.Errorf("failed to create directory: %w", err)
			}
			
			outFile, err := os.Create(targetPath)
			if err != nil {
				os.RemoveAll(tempDir) // Clean up on error
				return nil, fmt.Errorf("failed to create file: %w", err)
			}
			
			// Set file permissions
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				outFile.Close()
				os.RemoveAll(tempDir) // Clean up on error
				return nil, fmt.Errorf("failed to set file permissions: %w", err)
			}
			
			// Copy content
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				os.RemoveAll(tempDir) // Clean up on error
				return nil, fmt.Errorf("failed to copy file content: %w", err)
			}
			outFile.Close()
		}
	}
	
	// Load the runtime configuration
	config, err := LoadRuntimeConfig(tempDir)
	if err != nil {
		os.RemoveAll(tempDir) // Clean up on error
		return nil, fmt.Errorf("failed to load runtime configuration: %w", err)
	}
	
	// Get the entrypoint for the current OS
	entrypoint := GetEntrypointForCurrentOS()
	if entrypoint == "" {
		os.RemoveAll(tempDir) // Clean up on error
		return nil, fmt.Errorf("no entrypoint found for OS: %s", runtime.GOOS)
	}
	
	// Construct the executable path
	executablePath := filepath.Join(tempDir, entrypoint)
	
	// Ensure the executable has execute permissions
	if runtime.GOOS != "windows" {
		if err := os.Chmod(executablePath, 0755); err != nil {
			os.RemoveAll(tempDir) // Clean up on error
			return nil, fmt.Errorf("failed to set executable permissions: %w", err)
		}
	}
	
	return &Package{
		TempDir:    	tempDir,
		Metadata: 		config,
		Runtimepath: 	executablePath,
	}, nil
}

// CleanupWatcherPackage removes the temporary directory of an extracted package
func CleanupWatcherPackage(wp *Package) error {
	if wp != nil && wp.TempDir != "" {
		return os.RemoveAll(wp.TempDir)
	}
	return nil
}

// GetEntrypointForCurrentOS returns the entrypoint for the current OS
func GetEntrypointForCurrentOS() string {
	Targets := entrypoint.DefaultTargets
	var filename string
	
	for _, target := range Targets {
		if target.OS == runtime.GOOS && target.Arch == runtime.GOARCH {
			filename = target.Name
			break
		}
	}

	return filename
}