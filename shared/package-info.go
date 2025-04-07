package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pkginfo "rounds.com.ar/watcher/sdk/base/package-config"
)

// LoadRuntimeConfig loads package.info from the given directory
func LoadRuntimeConfig(dir string) (*pkginfo.PackageConfig, error) {
	infoPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.info: %w", err)
	}
	
	var config pkginfo.PackageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse package.info: %w", err)
	}
	
	return &config, nil
}
