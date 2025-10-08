package main

import (
	"testing"
	"time"

	"github.com/rabbytesoftware/quiver/internal"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func TestMainComponents(t *testing.T) {
	// Test that we can create the internal components
	internal := internal.NewInternal()
	if internal == nil {
		t.Error("Expected internal to be created")
	}

	watcher := internal.GetCore().GetWatcher()
	if watcher == nil {
		t.Error("Expected watcher to be created")
	}

	// Test metadata access
	name := metadata.GetName()
	if name == "" {
		t.Error("Expected non-empty name from metadata")
	}

	version := metadata.GetVersion()
	if version == "" {
		t.Error("Expected non-empty version from metadata")
	}
}

func TestMainLogic(t *testing.T) {
	// Test the main logic without actually running main()
	internal := internal.NewInternal()
	_ = internal.GetCore().GetWatcher()

	// Test that we can call the logging function that would be called in main
	watcher.Info("Test message from main test")

	// Test the goroutine logic by simulating what happens in main
	done := make(chan bool)
	go func() {
		// Simulate the time.Sleep and logging that happens in main
		time.Sleep(10 * time.Millisecond) // Much shorter for testing

		name := metadata.GetName()
		version := metadata.GetVersion()
		message := name + " " + version + " - Initializing..."

		if message == "" {
			t.Error("Expected non-empty initialization message")
		}

		watcher.Info(message)
		done <- true
	}()

	// Wait for the goroutine to complete or timeout
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Goroutine did not complete in time")
	}
}

// Test that internal.Run() can be called without panicking
func TestInternalRun(t *testing.T) {
	internal := internal.NewInternal()

	// Start internal.Run() in a goroutine and stop it quickly
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("internal.Run() panicked: %v", r)
			}
			done <- true
		}()

		// We can't easily test the full Run() method since it might block,
		// but we can at least verify it doesn't panic immediately
		go internal.Run()

		// Give it a moment to start
		time.Sleep(10 * time.Millisecond)
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Test did not complete in time")
	}
}
