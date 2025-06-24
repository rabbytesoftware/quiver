package packages

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"

	v0_1 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.1"
	yaml "gopkg.in/yaml.v3"
)

type ArrowsServer struct {
	Packages      map[string]manifest.ArrowInterface
	PackageManager *PackageManager
	
	logs *logger.Logger
}

func NewArrowsServer(repositories []string, installDir, dbPath string, logger *logger.Logger) *ArrowsServer {
	pm := NewPackageManager(repositories, installDir, dbPath, logger)
	
	return &ArrowsServer{
		Packages:       make(map[string]manifest.ArrowInterface),
		PackageManager: pm,
		logs:          logger, 
	}
}

// Initialize initializes the arrows server
func (as *ArrowsServer) Initialize() error {
	if err := as.PackageManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize package manager: %w", err)
	}

	// Load template arrows for backward compatibility
	for _, repo := range as.PackageManager.GetRepositories() {
		if err := as.LoadDirectory(repo); err != nil {
			as.logs.Warn("Failed to load templates from %s: %v", repo, err)
		}
	}

	return nil
}

func (as *ArrowsServer) Load(path string) error {
	// Read the YAML file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		as.logs.Error("Failed to read YAML file: %v", err)
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// First, parse just the version to determine which model to use
	var versionInfo struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal(data, &versionInfo); err != nil {
		as.logs.Error("Failed to parse version from YAML: %v", err)
		return fmt.Errorf("failed to parse version from YAML: %w", err)
	}

	if versionInfo.Version == "" {
		return fmt.Errorf("version field is required in YAML file")
	}

	// Load the package based on version
	var arrowObj manifest.ArrowInterface

	switch versionInfo.Version {
	case "0.1": // Create instance of v0.1 Arrow struct
		var arrow v0_1.Arrow
		if err := yaml.Unmarshal(data, &arrow); err != nil {
			as.logs.Error("Failed to unmarshal into v0.1 Arrow struct: %v", err)
			return fmt.Errorf("failed to unmarshal into v0.1 Arrow struct: %w", err)
		}
		arrowObj = &arrow
	default:
		return fmt.Errorf("unsupported version: %s", versionInfo.Version)
	}

	// Store the loaded package
	packageKey := fmt.Sprintf("%s @ %s ; with Arrow Manifest version: %s", 
		arrowObj.Name(), 
		arrowObj.ArrowVersion(), 
		versionInfo.Version)
	as.Packages[path] = arrowObj

	as.logs.Info("Successfully loaded package: %s from %s", packageKey, filepath.Base(path))

	return nil
}

func (as *ArrowsServer) GetPackage(path string) (manifest.ArrowInterface, bool) {
	pkg, exists := as.Packages[path]
	return pkg, exists
}

func (as *ArrowsServer) GetAllPackages() map[string]manifest.ArrowInterface {
	return as.Packages
}

func (as *ArrowsServer) LoadDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var loadErrors []error
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// Check if it's a YAML file
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			fullPath := filepath.Join(dirPath, file.Name())
			if err := as.Load(fullPath); err != nil {
				loadErrors = append(loadErrors, fmt.Errorf("failed to load %s: %w", file.Name(), err))
			}
		}
	}

	if len(loadErrors) > 0 {
		// Return the first error, but log all of them
		for i, err := range loadErrors {
			if i == 0 {
				as.logs.Error("Package loading errors encountered: %v", err)
			} else {
				as.logs.Error("Additional error: %v", err)
			}
		}
		return loadErrors[0]
	}

	return nil
}
