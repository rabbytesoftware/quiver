package commands

import (
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		input       string
		expectedCmd CommandKind
		expectError bool
		description string
	}{
		{"help", CmdHelp, false, "help command"},
		{"h", CmdHelp, false, "help shorthand"},
		{"HELP", CmdHelp, false, "help uppercase"},
		{"filter test", CmdFilter, false, "filter command with args"},
		{"f test", CmdFilter, false, "filter shorthand"},
		{"level info", CmdLevel, false, "level command"},
		{"l debug", CmdLevel, false, "level shorthand"},
		{"pause", CmdPause, false, "pause command"},
		{"p", CmdPause, false, "pause shorthand"},
		{"resume", CmdResume, false, "resume command"},
		{"r", CmdResume, false, "resume shorthand"},
		{"clear", CmdClear, false, "clear command"},
		{"c", CmdClear, false, "clear shorthand"},
		{"unknown command", CmdQuery, false, "unknown command becomes query"},
		{"", CmdQuery, true, "empty command"},
		{"   ", CmdQuery, true, "whitespace only command"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := Parse(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tc.input, err)
				return
			}

			if cmd.Kind != tc.expectedCmd {
				t.Errorf("Expected command kind %v, got %v for input %q", tc.expectedCmd, cmd.Kind, tc.input)
			}

			if cmd.OriginalInput != tc.input {
				t.Errorf("Expected OriginalInput %q, got %q", tc.input, cmd.OriginalInput)
			}
		})
	}
}

func TestParseWithArguments(t *testing.T) {
	testCases := []struct {
		input        string
		expectedArgs []string
		description  string
	}{
		{"help", []string{}, "help with no args"},
		{"filter test pattern", []string{"test", "pattern"}, "filter with multiple args"},
		{"level debug", []string{"debug"}, "level with one arg"},
		{"unknown command with args", []string{"command", "with", "args"}, "query with multiple args"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := Parse(tc.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(cmd.Args) != len(tc.expectedArgs) {
				t.Errorf("Expected %d args, got %d", len(tc.expectedArgs), len(cmd.Args))
				return
			}

			for i, expectedArg := range tc.expectedArgs {
				if cmd.Args[i] != expectedArg {
					t.Errorf("Expected arg[%d] = %q, got %q", i, expectedArg, cmd.Args[i])
				}
			}
		})
	}
}

func TestParseFilterRegexValidation(t *testing.T) {
	testCases := []struct {
		input       string
		expectError bool
		description string
	}{
		{"filter test", false, "valid simple pattern"},
		{"filter [a-z]+", false, "valid regex pattern"},
		{"filter .*", false, "valid wildcard pattern"},
		{"filter [", true, "invalid regex - unclosed bracket"},
		{"filter (", true, "invalid regex - unclosed paren"},
		{"filter \\", true, "invalid regex - trailing backslash"},
		{"f valid", false, "valid pattern with shorthand"},
		{"f [invalid", true, "invalid pattern with shorthand"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := Parse(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for invalid regex %q, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for valid regex %q: %v", tc.input, err)
				return
			}

			if cmd.Kind != CmdFilter {
				t.Errorf("Expected CmdFilter, got %v", cmd.Kind)
			}
		})
	}
}

func TestParseCommandValidation(t *testing.T) {
	// Test that Parse calls Validate() on the command
	validInputs := []string{
		"help",
		"filter test",
		"level info",
		"pause",
		"resume",
		"clear",
	}

	for _, input := range validInputs {
		t.Run("validate_"+input, func(t *testing.T) {
			cmd, err := Parse(input)

			if err != nil {
				t.Errorf("Valid command %q failed validation: %v", input, err)
				return
			}

			// Test that the command is valid by calling Validate directly
			if validateErr := cmd.Validate(); validateErr != nil {
				t.Errorf("Command validation failed for %q: %v", input, validateErr)
			}
		})
	}
}

func TestGetHelpText(t *testing.T) {
	helpText := GetHelpText()

	if helpText == "" {
		t.Error("GetHelpText() returned empty string")
	}

	// Test that help text contains expected commands
	expectedCommands := []string{
		"help",
		"filter",
		"level",
		"pause",
		"resume",
		"clear",
	}

	for _, cmd := range expectedCommands {
		if !contains(helpText, cmd) {
			t.Errorf("Help text does not contain command %q", cmd)
		}
	}

	// Test that help text contains navigation instructions
	expectedNavigation := []string{
		"Up/Down",
		"Ctrl+C",
		"Esc",
	}

	for _, nav := range expectedNavigation {
		if !contains(helpText, nav) {
			t.Errorf("Help text does not contain navigation %q", nav)
		}
	}
}

func TestParseEdgeCases(t *testing.T) {
	testCases := []struct {
		input       string
		expectError bool
		description string
	}{
		{"   help   ", false, "help with surrounding whitespace"},
		{"FILTER TEST", false, "uppercase filter"},
		{"Help", false, "mixed case help"},
		{"filter", true, "filter without args"},
		{"level", true, "level without args"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cmd, err := Parse(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %q, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %q: %v", tc.input, err)
				return
			}

			if cmd.OriginalInput != tc.input {
				t.Errorf("OriginalInput should preserve original input %q, got %q", tc.input, cmd.OriginalInput)
			}
		})
	}
}

func TestParseCommandKindMapping(t *testing.T) {
	// Test command mappings that don't require arguments
	mappings := map[string]CommandKind{
		"help":   CmdHelp,
		"h":      CmdHelp,
		"pause":  CmdPause,
		"p":      CmdPause,
		"resume": CmdResume,
		"r":      CmdResume,
		"clear":  CmdClear,
		"c":      CmdClear,
	}

	for input, expectedKind := range mappings {
		t.Run("mapping_"+input, func(t *testing.T) {
			cmd, err := Parse(input)

			if err != nil {
				t.Errorf("Error parsing %q: %v", input, err)
				return
			}

			if cmd.Kind != expectedKind {
				t.Errorf("Expected %v for %q, got %v", expectedKind, input, cmd.Kind)
			}
		})
	}

	// Test command mappings that require arguments
	argMappings := map[string]CommandKind{
		"filter test": CmdFilter,
		"f test":      CmdFilter,
		"level info":  CmdLevel,
		"l debug":     CmdLevel,
	}

	for input, expectedKind := range argMappings {
		t.Run("mapping_"+input, func(t *testing.T) {
			cmd, err := Parse(input)

			if err != nil {
				t.Errorf("Error parsing %q: %v", input, err)
				return
			}

			if cmd.Kind != expectedKind {
				t.Errorf("Expected %v for %q, got %v", expectedKind, input, cmd.Kind)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
