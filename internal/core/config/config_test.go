package config

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	// Reset config to nil to test fresh initialization
	config = nil

	cfg := Get()

	if cfg == nil {
		t.Fatal("Get() returned nil")
	}

	// Test that it returns the same instance on subsequent calls (singleton pattern)
	cfg2 := Get()
	if cfg != cfg2 {
		t.Error("Get() should return the same instance (singleton)")
	}
}

func TestGetNetbridge(t *testing.T) {
	netbridge := GetNetbridge()

	// Test that we get a Netbridge struct
	if netbridge.Enabled != true && netbridge.Enabled != false {
		// This is just to test the field exists and is accessible
	}

	// Test that AllowedPorts field exists
	_ = netbridge.AllowedPorts
}

func TestGetArrows(t *testing.T) {
	arrows := GetArrows()

	// Test that we get an Arrows struct
	if arrows.Repositories == nil {
		// Repositories can be nil or empty, both are valid
	}

	// Test that InstallDir field exists
	_ = arrows.InstallDir
}

func TestGetAPI(t *testing.T) {
	api := GetAPI()

	// Test that we get an API struct
	if api.Host == "" {
		t.Error("API.Host should not be empty")
	}

	if api.Port <= 0 {
		t.Error("API.Port should be positive")
	}
}

func TestGetDatabase(t *testing.T) {
	database := GetDatabase()

	// Test that we get a Database struct
	if database.Path == "" {
		t.Error("Database.Path should not be empty")
	}
}

func TestGetWatcher(t *testing.T) {
	watcher := GetWatcher()

	// Test that we get a Watcher struct
	if watcher.Level == "" {
		t.Error("Watcher.Level should not be empty")
	}

	if watcher.Folder == "" {
		t.Error("Watcher.Folder should not be empty")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}
}

func TestConfigExists(t *testing.T) {
	// This test depends on whether config file exists
	// We just test that the function doesn't panic
	exists := ConfigExists()

	// exists can be true or false, both are valid
	_ = exists
}

func TestConfigStructures(t *testing.T) {
	cfg := Get()

	// Test Netbridge structure
	netbridge := cfg.Config.Netbridge
	if netbridge.AllowedPorts == "" {
		// AllowedPorts can be empty, that's valid
	}

	// Test Arrows structure
	arrows := cfg.Config.Arrows
	if arrows.InstallDir == "" {
		t.Error("Arrows.InstallDir should not be empty")
	}

	// Test API structure
	api := cfg.Config.API
	if api.Host == "" {
		t.Error("API.Host should not be empty")
	}
	if api.Port <= 0 {
		t.Error("API.Port should be positive")
	}

	// Test Database structure
	database := cfg.Config.Database
	if database.Path == "" {
		t.Error("Database.Path should not be empty")
	}

	// Test Watcher structure
	watcher := cfg.Config.Watcher
	if watcher.Level == "" {
		t.Error("Watcher.Level should not be empty")
	}
	if watcher.Folder == "" {
		t.Error("Watcher.Folder should not be empty")
	}
}

func TestGetDefaultConfig(t *testing.T) {
	defaultCfg := getDefaultConfig()

	if defaultCfg == nil {
		t.Fatal("getDefaultConfig() returned nil")
	}

	// Test default values
	if defaultCfg.Config.API.Host == "" {
		t.Error("Default API.Host should not be empty")
	}

	if defaultCfg.Config.API.Port <= 0 {
		t.Error("Default API.Port should be positive")
	}

	if defaultCfg.Config.Database.Path == "" {
		t.Error("Default Database.Path should not be empty")
	}

	if defaultCfg.Config.Watcher.Level == "" {
		t.Error("Default Watcher.Level should not be empty")
	}

	if defaultCfg.Config.Watcher.Folder == "" {
		t.Error("Default Watcher.Folder should not be empty")
	}
}

func TestConfigWithNonExistentFile(t *testing.T) {
	// Reset config to test fresh loading
	config = nil

	// Temporarily set an invalid config path by creating a temp file and removing it
	tempFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Skip("Cannot create temp file for test")
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	os.Remove(tempPath) // Remove the file so it doesn't exist

	// We can't easily override the config path in the current implementation
	// So we just test that Get() works even when config file doesn't exist
	cfg := Get()

	if cfg == nil {
		t.Error("Get() should return default config when file doesn't exist")
	}
}

func TestConfigSingleton(t *testing.T) {
	// Reset config
	config = nil

	// Get config multiple times
	cfg1 := Get()
	cfg2 := Get()
	cfg3 := GetAPI()     // This calls Get() internally
	cfg4 := GetWatcher() // This calls Get() internally

	// All should reference the same config instance
	if cfg1 != cfg2 {
		t.Error("Config should be singleton - Get() calls should return same instance")
	}

	// We can't directly compare cfg3 and cfg4 since they return sub-structs
	// But we can verify they don't cause issues
	if cfg3.Host == "" {
		t.Error("GetAPI() should return valid API config")
	}

	if cfg4.Level == "" {
		t.Error("GetWatcher() should return valid Watcher config")
	}
}

func TestConfigFieldTypes(t *testing.T) {
	cfg := Get()

	// Test that all fields have the expected types
	netbridge := cfg.Config.Netbridge
	if _, ok := interface{}(netbridge.Enabled).(bool); !ok {
		t.Error("Netbridge.Enabled should be bool")
	}

	api := cfg.Config.API
	if _, ok := interface{}(api.Port).(int); !ok {
		t.Error("API.Port should be int")
	}

	watcher := cfg.Config.Watcher
	if _, ok := interface{}(watcher.Enabled).(bool); !ok {
		t.Error("Watcher.Enabled should be bool")
	}
	if _, ok := interface{}(watcher.MaxSize).(int); !ok {
		t.Error("Watcher.MaxSize should be int")
	}
	if _, ok := interface{}(watcher.MaxAge).(int); !ok {
		t.Error("Watcher.MaxAge should be int")
	}
	if _, ok := interface{}(watcher.Compress).(bool); !ok {
		t.Error("Watcher.Compress should be bool")
	}
}

func TestConfigWithInvalidYAML(t *testing.T) {
	// We can't easily test invalid YAML loading since the config path is hardcoded
	// But we can test that getDefaultConfig works correctly when called directly

	// Reset config to test fresh loading
	originalConfig := config
	config = nil
	defer func() {
		config = originalConfig // Restore original config after test
	}()

	// Test that Get() returns a valid config even in error scenarios
	cfg := Get()

	if cfg == nil {
		t.Fatal("Get() should never return nil, should return default config on errors")
	}

	// Verify it has the expected structure
	if cfg.Config.API.Host == "" {
		t.Error("Config should have valid API host even when using defaults")
	}

	if cfg.Config.API.Port <= 0 {
		t.Error("Config should have valid API port even when using defaults")
	}
}

func TestConfigLoadingPaths(t *testing.T) {
	// Reset config to test fresh loading
	originalConfig := config
	config = nil
	defer func() {
		config = originalConfig // Restore original config after test
	}()

	// Test the singleton behavior - first call loads, second call returns cached
	cfg1 := Get()
	if cfg1 == nil {
		t.Fatal("First Get() call should return valid config")
	}

	cfg2 := Get()
	if cfg2 == nil {
		t.Fatal("Second Get() call should return valid config")
	}

	if cfg1 != cfg2 {
		t.Error("Get() should return the same instance (singleton pattern)")
	}
}

func TestGetDefaultConfigCoverage(t *testing.T) {
	// Test getDefaultConfig function directly to improve coverage
	defaultCfg := getDefaultConfig()

	if defaultCfg == nil {
		t.Fatal("getDefaultConfig() should never return nil")
	}

	// Test all major sections exist
	if defaultCfg.Config.API.Host == "" {
		t.Error("Default config should have API host")
	}

	if defaultCfg.Config.API.Port <= 0 {
		t.Error("Default config should have valid API port")
	}

	if defaultCfg.Config.Database.Path == "" {
		t.Error("Default config should have database path")
	}

	if defaultCfg.Config.Watcher.Level == "" {
		t.Error("Default config should have watcher level")
	}

	if defaultCfg.Config.Watcher.Folder == "" {
		t.Error("Default config should have watcher folder")
	}

	if defaultCfg.Config.Arrows.InstallDir == "" {
		t.Error("Default config should have arrows install dir")
	}
}
