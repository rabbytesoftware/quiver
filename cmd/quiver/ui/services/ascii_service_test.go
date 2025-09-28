package services

import (
	"testing"
)

func TestNewASCIIService(t *testing.T) {
	service := NewASCIIService()
	if service == nil {
		t.Fatal("Expected ASCII service to be created")
	}
}

func TestASCIIService_GetWelcomeLogLine(t *testing.T) {
	service := NewASCIIService()

	logLine := service.GetWelcomeLogLine()

	if logLine.Text == "" {
		t.Error("Expected ASCII art text to be present")
	}

	if logLine.Level != "info" {
		t.Errorf("Expected level to be 'info', got '%s'", logLine.Level)
	}

	if logLine.Time.IsZero() {
		t.Error("Expected time to be set")
	}

	// Verify the ASCII art contains the expected pattern
	if len(logLine.Text) < 100 {
		t.Error("Expected ASCII art to be substantial")
	}
}
