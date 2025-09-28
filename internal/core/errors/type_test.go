package errors

import (
	"testing"
)

func TestNewErr(t *testing.T) {
	testCases := []struct {
		name    string
		code    int
		message string
	}{
		{
			name:    "basic error",
			code:    404,
			message: "Not found",
		},
		{
			name:    "server error",
			code:    500,
			message: "Internal server error",
		},
		{
			name:    "zero code",
			code:    0,
			message: "Unknown error",
		},
		{
			name:    "negative code",
			code:    -1,
			message: "Invalid error code",
		},
		{
			name:    "empty message",
			code:    400,
			message: "",
		},
		{
			name:    "long message",
			code:    422,
			message: "This is a very long error message that describes in detail what went wrong with the request processing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewErr(tc.code, tc.message)

			if err == nil {
				t.Fatal("NewErr() returned nil")
			}

			if err.Code != tc.code {
				t.Errorf("Expected Code to be %d, got %d", tc.code, err.Code)
			}

			if err.Message != tc.message {
				t.Errorf("Expected Message to be %q, got %q", tc.message, err.Message)
			}
		})
	}
}

func TestErr_Error(t *testing.T) {
	testCases := []struct {
		name          string
		code          int
		message       string
		expectedError string
	}{
		{
			name:          "basic error",
			code:          404,
			message:       "Not found",
			expectedError: "Not found",
		},
		{
			name:          "empty message",
			code:          500,
			message:       "",
			expectedError: "",
		},
		{
			name:          "message with special characters",
			code:          400,
			message:       "Error: invalid input (code: 123)",
			expectedError: "Error: invalid input (code: 123)",
		},
		{
			name:          "unicode message",
			code:          422,
			message:       "エラーが発生しました",
			expectedError: "エラーが発生しました",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewErr(tc.code, tc.message)
			result := err.Error()

			if result == nil {
				t.Fatal("Error() returned nil")
			}

			if result.Error() != tc.expectedError {
				t.Errorf("Expected Error() to return %q, got %q", tc.expectedError, result.Error())
			}
		})
	}
}

func TestErr_StructFields(t *testing.T) {
	code := 418
	message := "I'm a teapot"

	err := NewErr(code, message)

	// Test direct field access
	if err.Code != code {
		t.Errorf("Expected Code field to be %d, got %d", code, err.Code)
	}

	if err.Message != message {
		t.Errorf("Expected Message field to be %q, got %q", message, err.Message)
	}

	// Test that fields can be modified
	err.Code = 500
	err.Message = "Modified message"

	if err.Code != 500 {
		t.Errorf("Expected modified Code to be 500, got %d", err.Code)
	}

	if err.Message != "Modified message" {
		t.Errorf("Expected modified Message to be 'Modified message', got %q", err.Message)
	}
}

func TestErr_ZeroValue(t *testing.T) {
	var err Err

	// Test zero values
	if err.Code != 0 {
		t.Errorf("Expected zero value Code to be 0, got %d", err.Code)
	}

	if err.Message != "" {
		t.Errorf("Expected zero value Message to be empty, got %q", err.Message)
	}

	// Test Error() method on zero value
	result := err.Error()
	if result == nil {
		t.Fatal("Error() on zero value returned nil")
	}

	if result.Error() != "" {
		t.Errorf("Expected Error() on zero value to return empty string, got %q", result.Error())
	}
}

func TestErr_MultipleInstances(t *testing.T) {
	err1 := NewErr(404, "Not found")
	err2 := NewErr(500, "Server error")
	err3 := NewErr(404, "Not found")

	// Test that different instances are different objects
	if err1 == err2 {
		t.Error("Different error instances should not be equal")
	}

	if err1 == err3 {
		t.Error("Different error instances should not be equal even with same values")
	}

	// Test that values are independent
	if err1.Code == err2.Code {
		t.Error("Different errors should have different codes")
	}

	if err1.Message == err2.Message {
		t.Error("Different errors should have different messages")
	}

	// Test that same values produce same results
	if err1.Code != err3.Code {
		t.Error("Errors with same code should have same Code field")
	}

	if err1.Message != err3.Message {
		t.Error("Errors with same message should have same Message field")
	}
}

func TestErr_ErrorInterface(t *testing.T) {
	err := NewErr(400, "Bad request")

	// Test that Err implements the error interface pattern
	// (though it returns error, not implements error directly)
	result := err.Error()

	// Verify it returns an error type
	if result == nil {
		t.Fatal("Error() should return an error, got nil")
	}

	// Test that the error message is correct
	expectedMessage := "Bad request"
	if result.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, result.Error())
	}
}

func TestErr_EdgeCases(t *testing.T) {
	// Test with maximum int value
	maxErr := NewErr(int(^uint(0)>>1), "Max int error")
	if maxErr.Code != int(^uint(0)>>1) {
		t.Error("Should handle maximum int value")
	}

	// Test with minimum int value
	minErr := NewErr(-int(^uint(0)>>1)-1, "Min int error")
	if minErr.Code != -int(^uint(0)>>1)-1 {
		t.Error("Should handle minimum int value")
	}

	// Test with very long message
	longMessage := make([]byte, 10000)
	for i := range longMessage {
		longMessage[i] = 'a'
	}
	longErr := NewErr(200, string(longMessage))
	if len(longErr.Message) != 10000 {
		t.Error("Should handle very long messages")
	}

	// Test that Error() works with long message
	result := longErr.Error()
	if result == nil {
		t.Error("Error() should work with long messages")
	}
}

func TestErr_Consistency(t *testing.T) {
	err := NewErr(403, "Forbidden")

	// Test that multiple calls to Error() return consistent results
	result1 := err.Error()
	result2 := err.Error()

	if result1.Error() != result2.Error() {
		t.Error("Multiple calls to Error() should return consistent results")
	}

	// Test that modifying the struct affects Error() output
	err.Message = "Modified forbidden"
	result3 := err.Error()

	if result3.Error() != "Modified forbidden" {
		t.Error("Error() should reflect changes to Message field")
	}
}
