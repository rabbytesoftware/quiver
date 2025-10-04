package assets

import (
	"os"
	"runtime"
	"testing"
)

func TestIconManager_GetIconForPlatform(t *testing.T) {
	manager := NewIconManager()
	iconData := manager.GetIconForPlatform()

	if len(iconData) == 0 {
		t.Fatal("Icon data should not be empty")
	}

	// Test platform-specific icon selection
	switch runtime.GOOS {
	case "windows":
		if len(WindowsIcon) == 0 {
			t.Error("Windows icon should be embedded")
		}
	case "darwin":
		if len(MacOSIcon) == 0 {
			t.Error("macOS icon should be embedded")
		}
	case "linux":
		if len(LinuxIcon) == 0 {
			t.Error("Linux icon should be embedded")
		}
	}
}

func TestIconManager_SaveIconToTemp(t *testing.T) {
	manager := NewIconManager()
	tempFile, err := manager.SaveIconToTemp()
	if err != nil {
		t.Fatalf("Failed to save icon to temp: %v", err)
	}
	defer manager.Cleanup()

	if tempFile == "" {
		t.Error("Temp file path should not be empty")
	}

	// Check if file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Temp icon file should exist")
	}

	// Check file size
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Icon file should not be empty")
	}
}

func TestIconManager_GetIconSize(t *testing.T) {
	manager := NewIconManager()
	iconData := manager.GetIconSize()

	if len(iconData) == 0 {
		t.Error("Icon size data should not be empty")
	}
}

func TestIconManager_GetAvailableSizes(t *testing.T) {
	manager := NewIconManager()
	sizes := manager.GetAvailableSizes()

	expectedSizes := []string{"16", "32", "64", "128", "256"}
	for _, size := range expectedSizes {
		if _, exists := sizes[size]; !exists {
			t.Errorf("Size %s should be available", size)
		}
		if len(sizes[size]) == 0 {
			t.Errorf("Size %s should have data", size)
		}
	}
}

func TestEmbeddedIcons(t *testing.T) {
	// Test that all embedded icons have data
	icons := map[string][]byte{
		"WindowsIcon": WindowsIcon,
		"MacOSIcon":   MacOSIcon,
		"LinuxIcon":   LinuxIcon,
		"Icon16":      Icon16,
		"Icon32":      Icon32,
		"Icon64":      Icon64,
		"Icon128":     Icon128,
		"Icon256":     Icon256,
	}

	for name, data := range icons {
		if len(data) == 0 {
			t.Errorf("Embedded icon %s should have data", name)
		}
	}
}

func TestIconManager_GetIconPath(t *testing.T) {
	manager := NewIconManager()
	
	// Initially should be empty
	if manager.GetIconPath() != "" {
		t.Error("Icon path should be empty initially")
	}
	
	// Save icon to temp and check path
	tempFile, err := manager.SaveIconToTemp()
	if err != nil {
		t.Fatalf("Failed to save icon to temp: %v", err)
	}
	defer manager.Cleanup()
	
	if manager.GetIconPath() != tempFile {
		t.Error("Icon path should match saved temp file")
	}
}

func TestIconManager_Cleanup(t *testing.T) {
	manager := NewIconManager()
	
	// Test cleanup when no icon is saved
	err := manager.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error when no icon is saved: %v", err)
	}
	
	// Save icon and test cleanup
	tempFile, err := manager.SaveIconToTemp()
	if err != nil {
		t.Fatalf("Failed to save icon to temp: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Icon file should exist before cleanup")
	}
	
	// Cleanup and verify file is removed
	err = manager.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error: %v", err)
	}
	
	// Verify file is removed
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("Icon file should be removed after cleanup")
	}
}

func TestIconManager_PlatformSpecific(t *testing.T) {
	manager := NewIconManager()
	
	// Test GetIconForPlatform for all platforms
	iconData := manager.GetIconForPlatform()
	if len(iconData) == 0 {
		t.Error("Icon data should not be empty for any platform")
	}
	
	// Test GetIconSize for all platforms
	sizeData := manager.GetIconSize()
	if len(sizeData) == 0 {
		t.Error("Icon size data should not be empty for any platform")
	}
	
	// Test platform-specific behavior
	switch runtime.GOOS {
	case "windows":
		// Windows should return WindowsIcon
		if len(WindowsIcon) == 0 {
			t.Error("Windows icon should be embedded")
		}
		// Windows should use 32x32 icon size
		if len(Icon32) == 0 {
			t.Error("32x32 icon should be embedded")
		}
	case "darwin":
		// macOS should return MacOSIcon
		if len(MacOSIcon) == 0 {
			t.Error("macOS icon should be embedded")
		}
		// macOS should use 256x256 icon size
		if len(Icon256) == 0 {
			t.Error("256x256 icon should be embedded")
		}
	case "linux":
		// Linux should return LinuxIcon
		if len(LinuxIcon) == 0 {
			t.Error("Linux icon should be embedded")
		}
		// Linux should use 128x128 icon size
		if len(Icon128) == 0 {
			t.Error("128x128 icon should be embedded")
		}
	default:
		// Unknown platform should default to Linux behavior
		if len(LinuxIcon) == 0 {
			t.Error("Linux icon should be embedded for unknown platforms")
		}
		if len(Icon128) == 0 {
			t.Error("128x128 icon should be embedded for unknown platforms")
		}
	}
}

func TestIconManager_MultipleOperations(t *testing.T) {
	manager := NewIconManager()
	
	// Test multiple save operations
	tempFile1, err := manager.SaveIconToTemp()
	if err != nil {
		t.Fatalf("First save failed: %v", err)
	}
	defer manager.Cleanup()
	
	// Save again (should overwrite)
	tempFile2, err := manager.SaveIconToTemp()
	if err != nil {
		t.Fatalf("Second save failed: %v", err)
	}
	
	// Paths should be the same
	if tempFile1 != tempFile2 {
		t.Error("Icon path should remain the same on multiple saves")
	}
	
	// Test that file exists and has content
	if _, err := os.Stat(tempFile2); os.IsNotExist(err) {
		t.Error("Icon file should exist after second save")
	}
}
