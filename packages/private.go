package packages

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	entrypoint "github.com/rabbytesoftware/quiver.compiler/shared/base/entrypoint"
	pc "github.com/rabbytesoftware/quiver.compiler/shared/base/package-config"
	pb "github.com/rabbytesoftware/quiver.compiler/shared/package"
)

func (pkg *Package) exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
			
	_, err := pkg.Runtime.Client.Exit(ctx, &pb.Empty{})
	return err
}

func (pkg *Package) extract() error {
	// Create a unique temporary directory
	pkg.Runtime.TempDir = filepath.Join(os.TempDir() + "/watcher/", fmt.Sprintf("%s-%d", filepath.Base(pkg.Runtime.PackagePath), os.Getpid()))
	
	if err := os.MkdirAll(pkg.Runtime.TempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	// Open the package file
	file, err := os.Open(pkg.Runtime.PackagePath)
	if err != nil {
		os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
		return fmt.Errorf("failed to open package file: %w", err)
	}
	defer file.Close()
	
	// Create gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		os.RemoveAll(pkg.Runtime.PackagePath) // Clean up on error
		return fmt.Errorf("failed to create gzip reader: %w", err)
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
			os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
			return fmt.Errorf("error reading archive: %w", err)
		}
		
		targetPath := filepath.Join(pkg.Runtime.TempDir, header.Name)
		
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create file
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
				return fmt.Errorf("failed to create directory: %w", err)
			}
			
			outFile, err := os.Create(targetPath)
			if err != nil {
				os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
				return fmt.Errorf("failed to create file: %w", err)
			}
			
			// Set file permissions
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				outFile.Close()
				os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
				return fmt.Errorf("failed to set file permissions: %w", err)
			}
			
			// Copy content
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
				return fmt.Errorf("failed to copy file content: %w", err)
			}
			outFile.Close()
		}
	}
	
	// Load the runtime configuration
	err = pkg.loadRuntimeConfig()
	if err != nil {
		os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
		return fmt.Errorf("failed to load runtime configuration: %w", err)
	}
	
	// Get the entrypoint for the current OS
	entrypoint := getEntrypointForCurrentOS()
	if entrypoint == "" {
		os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
		return fmt.Errorf("no entrypoint found for OS: %s", runtime.GOOS)
	}
	
	// Construct the executable path
	pkg.Runtime.RuntimePath = filepath.Join(pkg.Runtime.TempDir, entrypoint)
	
	// Ensure the executable has execute permissions
	if runtime.GOOS != "windows" {
		if err := os.Chmod(pkg.Runtime.RuntimePath, 0755); err != nil {
			os.RemoveAll(pkg.Runtime.TempDir) // Clean up on error
			return fmt.Errorf("failed to set executable permissions: %w", err)
		}
	}
	
	return nil
}

func (pkg *Package) clean() error {
	if pkg.Runtime.TempDir != "" {
		return os.RemoveAll(pkg.Runtime.TempDir)
	}

	return nil
}

func (pkg *Package) loadRuntimeConfig() error {
	infoPath := filepath.Join(pkg.Runtime.TempDir, "package.json")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return fmt.Errorf("failed to read package.info: %w", err)
	}
	
	var config pc.PackageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse package.info: %w", err)
	}
	
	pkg.PackageConfig = &config

	return nil
}

func getEntrypointForCurrentOS() string {
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
