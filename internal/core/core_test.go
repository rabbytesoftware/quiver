package core

import (
	"testing"
)

func TestInit(t *testing.T) {
	core := Init()

	if core == nil {
		t.Fatal("Init() returned nil")
	}

	// Test that all dependencies are initialized
	if core.metadata == nil {
		t.Error("metadata is not initialized")
	}

	if core.config == nil {
		t.Error("config is not initialized")
	}

	if core.watcher == nil {
		t.Error("watcher is not initialized")
	}
}

func TestCore_GetMetadata(t *testing.T) {
	core := Init()
	metadata := core.GetMetadata()

	if metadata == nil {
		t.Error("GetMetadata() returned nil")
	}

	// Test that it returns the same instance
	metadata2 := core.GetMetadata()
	if metadata != metadata2 {
		t.Error("GetMetadata() should return the same instance")
	}

	// Test that metadata is valid
	if metadata == nil {
		t.Error("metadata is nil")
	}
}

func TestCore_GetConfig(t *testing.T) {
	core := Init()
	config := core.GetConfig()

	if config == nil {
		t.Error("GetConfig() returned nil")
	}

	// Test that it returns the same instance
	config2 := core.GetConfig()
	if config != config2 {
		t.Error("GetConfig() should return the same instance")
	}

	// Test that config is valid
	if config == nil {
		t.Error("config is nil")
	}
}

func TestCore_GetWatcher(t *testing.T) {
	core := Init()
	watcher := core.GetWatcher()

	if watcher == nil {
		t.Error("GetWatcher() returned nil")
	}

	// Test that it returns the same instance
	watcher2 := core.GetWatcher()
	if watcher != watcher2 {
		t.Error("GetWatcher() should return the same instance")
	}

	// Test that watcher is valid
	if watcher == nil {
		t.Error("watcher is nil")
	}
}

func TestCoreStructure(t *testing.T) {
	core := Init()

	// Test that Core struct has the expected fields
	if core.metadata == nil {
		t.Error("Core.metadata field is nil")
	}

	if core.config == nil {
		t.Error("Core.config field is nil")
	}

	if core.watcher == nil {
		t.Error("Core.watcher field is nil")
	}
}

func TestCoreInitialization(t *testing.T) {
	// Test multiple initializations to ensure consistency
	core1 := Init()
	core2 := Init()

	// Each Init() call should create a new Core instance
	if core1 == core2 {
		t.Error("Init() should create new instances each time")
	}

	// But the underlying services should be singletons
	if core1.GetMetadata() != core2.GetMetadata() {
		t.Error("Metadata should be singleton across Core instances")
	}

	if core1.GetConfig() != core2.GetConfig() {
		t.Error("Config should be singleton across Core instances")
	}

	// Watcher uses singleton pattern, so they should be the same
	if core1.GetWatcher() != core2.GetWatcher() {
		t.Error("Watcher should be singleton across Core instances")
	}
}
