package assets

import (
	"os"
	"path/filepath"
	"runtime"
)

// IconManager handles icon operations for different platforms
type IconManager struct {
	iconPath string
}

// NewIconManager creates a new IconManager instance
func NewIconManager() *IconManager {
	return &IconManager{}
}

// GetIconForPlatform returns the appropriate icon data for the current platform
func (im *IconManager) GetIconForPlatform() []byte {
	switch runtime.GOOS {
	case "windows":
		return WindowsIcon
	case "darwin":
		return MacOSIcon
	case "linux":
		return LinuxIcon
	default:
		return LinuxIcon // Default to PNG for unknown platforms
	}
}

// SaveIconToTemp saves the platform-appropriate icon to a temporary file
func (im *IconManager) SaveIconToTemp() (string, error) {
	iconData := im.GetIconForPlatform()

	// Create temp directory if it doesn't exist
	tempDir := os.TempDir()
	iconDir := filepath.Join(tempDir, "quiver-icons")
	if err := os.MkdirAll(iconDir, 0755); err != nil {
		return "", err
	}

	// Determine file extension based on platform
	var ext string
	switch runtime.GOOS {
	case "windows":
		ext = ".ico"
	case "darwin":
		ext = ".icns"
	case "linux":
		ext = ".png"
	default:
		ext = ".png"
	}

	// Create temp file
	tempFile := filepath.Join(iconDir, "quiver-icon"+ext)
	if err := os.WriteFile(tempFile, iconData, 0644); err != nil {
		return "", err
	}

	im.iconPath = tempFile
	return tempFile, nil
}

// GetIconPath returns the current icon path
func (im *IconManager) GetIconPath() string {
	return im.iconPath
}

// Cleanup removes the temporary icon file
func (im *IconManager) Cleanup() error {
	if im.iconPath != "" {
		return os.Remove(im.iconPath)
	}
	return nil
}

// GetIconSize returns the appropriate icon size for the current platform
func (im *IconManager) GetIconSize() []byte {
	switch runtime.GOOS {
	case "windows":
		return Icon32 // Windows typically uses 32x32
	case "darwin":
		return Icon256 // macOS typically uses 256x256
	case "linux":
		return Icon128 // Linux typically uses 128x128
	default:
		return Icon128
	}
}

// GetAvailableSizes returns all available icon sizes
func (im *IconManager) GetAvailableSizes() map[string][]byte {
	return map[string][]byte{
		"16":  Icon16,
		"32":  Icon32,
		"64":  Icon64,
		"128": Icon128,
		"256": Icon256,
	}
}
