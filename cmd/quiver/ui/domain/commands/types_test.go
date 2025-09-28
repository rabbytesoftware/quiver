package commands

import (
	"testing"
)

func TestCommandKind_String(t *testing.T) {
	testCases := []struct {
		kind     CommandKind
		expected string
	}{
		{CmdHelp, "help"},
		{CmdFilter, "filter"},
		{CmdLevel, "level"},
		{CmdPause, "pause"},
		{CmdResume, "resume"},
		{CmdClear, "clear"},
		{CmdQuery, "query"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.kind.String()
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestCommand_Validate(t *testing.T) {
	testCases := []struct {
		cmd         Command
		expectError bool
		description string
	}{
		{
			Command{Kind: CmdHelp, Args: []string{}, OriginalInput: "help"},
			false,
			"valid help command",
		},
		{
			Command{Kind: CmdHelp, Args: []string{"extra"}, OriginalInput: "help extra"},
			true,
			"help with args should be invalid",
		},
		{
			Command{Kind: CmdFilter, Args: []string{"pattern"}, OriginalInput: "filter pattern"},
			false,
			"valid filter command",
		},
		{
			Command{Kind: CmdFilter, Args: []string{}, OriginalInput: "filter"},
			true,
			"filter without args should be invalid",
		},
		{
			Command{Kind: CmdLevel, Args: []string{"info"}, OriginalInput: "level info"},
			false,
			"valid level command",
		},
		{
			Command{Kind: CmdLevel, Args: []string{}, OriginalInput: "level"},
			true,
			"level without args should be invalid",
		},
		{
			Command{Kind: CmdLevel, Args: []string{"invalid"}, OriginalInput: "level invalid"},
			true,
			"level with invalid level should be invalid",
		},
		{
			Command{Kind: CmdPause, Args: []string{}, OriginalInput: "pause"},
			false,
			"valid pause command",
		},
		{
			Command{Kind: CmdPause, Args: []string{"extra"}, OriginalInput: "pause extra"},
			true,
			"pause with args should be invalid",
		},
		{
			Command{Kind: CmdResume, Args: []string{}, OriginalInput: "resume"},
			false,
			"valid resume command",
		},
		{
			Command{Kind: CmdResume, Args: []string{"extra"}, OriginalInput: "resume extra"},
			true,
			"resume with args should be invalid",
		},
		{
			Command{Kind: CmdClear, Args: []string{}, OriginalInput: "clear"},
			false,
			"valid clear command",
		},
		{
			Command{Kind: CmdClear, Args: []string{"extra"}, OriginalInput: "clear extra"},
			true,
			"clear with args should be invalid",
		},
		{
			Command{Kind: CmdQuery, Args: []string{"any", "args"}, OriginalInput: "query any args"},
			false,
			"query commands are always valid",
		},
		{
			Command{Kind: CmdQuery, Args: []string{}, OriginalInput: "query"},
			false,
			"query without args is valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.cmd.Validate()

			if tc.expectError && err == nil {
				t.Error("Expected validation error, but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestValidLogLevels(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}

	for _, level := range validLevels {
		t.Run("level_"+level, func(t *testing.T) {
			cmd := Command{
				Kind:          CmdLevel,
				Args:          []string{level},
				OriginalInput: "level " + level,
			}

			err := cmd.Validate()
			if err != nil {
				t.Errorf("Level %q should be valid, got error: %v", level, err)
			}
		})
	}

	// Test invalid levels
	invalidLevels := []string{"trace", "fatal", "panic", "verbose", "invalid"}

	for _, level := range invalidLevels {
		t.Run("invalid_level_"+level, func(t *testing.T) {
			cmd := Command{
				Kind:          CmdLevel,
				Args:          []string{level},
				OriginalInput: "level " + level,
			}

			err := cmd.Validate()
			if err == nil {
				t.Errorf("Level %q should be invalid", level)
			}
		})
	}
}

func TestNewCommandError(t *testing.T) {
	message := "test error message"
	input := "test input"

	err := NewCommandError(message, input)

	if err == nil {
		t.Fatal("NewCommandError returned nil")
	}

	// Test that it implements error interface
	var _ error = err

	// Test error message
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Error message should not be empty")
	}

	// Test that error message contains the input message
	if !contains(errorMsg, message) {
		t.Errorf("Error message should contain %q, got %q", message, errorMsg)
	}
}

func TestCommandErrorType(t *testing.T) {
	err := NewCommandError("test", "input")

	// Test that it's a CommandError type
	if cmdErr, ok := err.(CommandError); !ok {
		t.Error("NewCommandError should return CommandError type")
	} else {
		// Test the fields are set correctly
		if cmdErr.Message != "test" {
			t.Error("CommandError.Message not set correctly")
		}
		if cmdErr.Input != "input" {
			t.Error("CommandError.Input not set correctly")
		}
	}

	// Test that it implements error interface (err is already error type)
	_ = err
}

func TestCommandStructure(t *testing.T) {
	cmd := Command{
		Kind:          CmdHelp,
		Args:          []string{"arg1", "arg2"},
		OriginalInput: "help arg1 arg2",
	}

	// Test field access
	if cmd.Kind != CmdHelp {
		t.Error("Command.Kind not set correctly")
	}

	if len(cmd.Args) != 2 {
		t.Error("Command.Args not set correctly")
	}

	if cmd.Args[0] != "arg1" || cmd.Args[1] != "arg2" {
		t.Error("Command.Args values not set correctly")
	}

	if cmd.OriginalInput != "help arg1 arg2" {
		t.Error("Command.OriginalInput not set correctly")
	}
}

func TestCommandKindConstants(t *testing.T) {
	// Test that all CommandKind constants are defined and unique
	kinds := []CommandKind{
		CmdHelp,
		CmdFilter,
		CmdLevel,
		CmdPause,
		CmdResume,
		CmdClear,
		CmdQuery,
	}

	// Test that they have different values
	seen := make(map[CommandKind]bool)
	for _, kind := range kinds {
		if seen[kind] {
			t.Errorf("Duplicate CommandKind value: %v", kind)
		}
		seen[kind] = true
	}

	// Test that they all have string representations
	for _, kind := range kinds {
		str := kind.String()
		if str == "" {
			t.Errorf("CommandKind %v has empty string representation", kind)
		}
	}
}
