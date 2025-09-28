package runtime

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/models/system"
)

func TestMethod_Structure(t *testing.T) {
	// Test that Method struct can be created and has expected fields
	method := Method{
		OS:      system.OS("linux/amd64"),
		Command: []string{"echo", "hello", "world"},
	}

	// Test field access
	if method.OS != system.OS("linux/amd64") {
		t.Errorf("Expected OS 'linux/amd64', got %q", method.OS)
	}

	if len(method.Command) != 3 {
		t.Errorf("Expected 3 command parts, got %d", len(method.Command))
	}

	if method.Command[0] != "echo" {
		t.Errorf("Expected first command part 'echo', got %q", method.Command[0])
	}

	if method.Command[1] != "hello" {
		t.Errorf("Expected second command part 'hello', got %q", method.Command[1])
	}

	if method.Command[2] != "world" {
		t.Errorf("Expected third command part 'world', got %q", method.Command[2])
	}
}

func TestMethod_EmptyMethod(t *testing.T) {
	// Test empty method
	method := Method{}

	if method.OS != "" {
		t.Errorf("Expected empty OS, got %q", method.OS)
	}

	if method.Command != nil {
		t.Errorf("Expected nil Command, got %v", method.Command)
	}
}

func TestMethod_OSTypes(t *testing.T) {
	// Test different OS types
	methods := []Method{
		{OS: system.OS("linux/amd64"), Command: []string{"ls", "-la"}},
		{OS: system.OS("windows/amd64"), Command: []string{"dir"}},
		{OS: system.OS("darwin/arm64"), Command: []string{"ls", "-la"}},
	}

	for i, method := range methods {
		if !method.OS.IsValid() {
			t.Errorf("Method %d: Expected OS to be valid, got %q", i, method.OS)
		}

		// Test OS type methods
		switch method.OS {
		case system.OS("linux/amd64"):
			if !method.OS.IsLinux() {
				t.Errorf("Method %d: Expected OS to be Linux", i)
			}
			if !method.OS.IsAMD64() {
				t.Errorf("Method %d: Expected OS to be AMD64", i)
			}
		case system.OS("windows/amd64"):
			if !method.OS.IsWindows() {
				t.Errorf("Method %d: Expected OS to be Windows", i)
			}
			if !method.OS.IsAMD64() {
				t.Errorf("Method %d: Expected OS to be AMD64", i)
			}
		case system.OS("darwin/arm64"):
			if !method.OS.IsDarwin() {
				t.Errorf("Method %d: Expected OS to be Darwin", i)
			}
			if !method.OS.IsARM64() {
				t.Errorf("Method %d: Expected OS to be ARM64", i)
			}
		}
	}
}

func TestMethod_CommandVariations(t *testing.T) {
	// Test different command variations
	testCases := []struct {
		name     string
		method   Method
		expected int
	}{
		{
			name:     "single command",
			method:   Method{Command: []string{"echo"}},
			expected: 1,
		},
		{
			name:     "command with arguments",
			method:   Method{Command: []string{"ls", "-la", "/tmp"}},
			expected: 3,
		},
		{
			name:     "complex command",
			method:   Method{Command: []string{"docker", "run", "--rm", "-it", "ubuntu:latest", "bash"}},
			expected: 6,
		},
		{
			name:     "empty command",
			method:   Method{Command: []string{}},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.method.Command) != tc.expected {
				t.Errorf("Expected %d command parts, got %d", tc.expected, len(tc.method.Command))
			}
		})
	}
}

func TestMethod_CombinedFields(t *testing.T) {
	// Test method with both OS and Command
	method := Method{
		OS:      system.OS("linux/amd64"),
		Command: []string{"docker", "ps", "-a"},
	}

	// Verify OS
	if !method.OS.IsValid() {
		t.Error("Expected OS to be valid")
	}

	if !method.OS.IsLinux() {
		t.Error("Expected OS to be Linux")
	}

	// Verify Command
	if len(method.Command) != 3 {
		t.Errorf("Expected 3 command parts, got %d", len(method.Command))
	}

	expectedCommand := []string{"docker", "ps", "-a"}
	for i, part := range method.Command {
		if part != expectedCommand[i] {
			t.Errorf("Expected command part %d to be %q, got %q", i, expectedCommand[i], part)
		}
	}
}
