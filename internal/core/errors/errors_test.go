package errors

import (
	"fmt"
	"testing"
)

func TestThrow(t *testing.T) {
	tests := []struct {
		name    string
		code    ErrorCode
		message string
		details map[string]interface{}
		want    Error
	}{
		{
			name:    "Basic error creation",
			code:    NotFound,
			message: "Resource not found",
			details: nil,
			want: Error{
				Code:    NotFound,
				Message: "Resource not found",
				Details: nil,
			},
		},
		{
			name:    "Error with details",
			code:    InvalidRequest,
			message: "Invalid input",
			details: map[string]interface{}{
				"field":  "email",
				"reason": "invalid format",
				"value":  "not-an-email",
			},
			want: Error{
				Code:    InvalidRequest,
				Message: "Invalid input",
				Details: map[string]interface{}{
					"field":  "email",
					"reason": "invalid format",
					"value":  "not-an-email",
				},
			},
		},
		{
			name:    "Server error",
			code:    InternalServer,
			message: "Database connection failed",
			details: map[string]interface{}{
				"retry_count": 3,
				"timeout":     "30s",
			},
			want: Error{
				Code:    InternalServer,
				Message: "Database connection failed",
				Details: map[string]interface{}{
					"retry_count": 3,
					"timeout":     "30s",
				},
			},
		},
		{
			name:    "Empty message",
			code:    Unauthorized,
			message: "",
			details: nil,
			want: Error{
				Code:    Unauthorized,
				Message: "",
				Details: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Throw(tt.code, tt.message, tt.details)

			if got.Code != tt.want.Code {
				t.Errorf("Throw() Code = %v, want %v", got.Code, tt.want.Code)
			}
			if got.Message != tt.want.Message {
				t.Errorf("Throw() Message = %v, want %v", got.Message, tt.want.Message)
			}
			if got.Details == nil && tt.want.Details != nil {
				t.Errorf("Throw() Details = %v, want %v", got.Details, tt.want.Details)
			}
			if got.Details != nil && tt.want.Details != nil {
				if len(got.Details) != len(tt.want.Details) {
					t.Errorf("Throw() Details length = %v, want %v", len(got.Details), len(tt.want.Details))
				}
				for k, v := range tt.want.Details {
					if got.Details[k] != v {
						t.Errorf("Throw() Details[%s] = %v, want %v", k, got.Details[k], v)
					}
				}
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "Basic error message",
			err: Error{
				Code:    NotFound,
				Message: "Resource not found",
			},
			want: "404: Resource not found",
		},
		{
			name: "Server error message",
			err: Error{
				Code:    InternalServer,
				Message: "Database connection failed",
			},
			want: "500: Database connection failed",
		},
		{
			name: "Empty message",
			err: Error{
				Code:    Unauthorized,
				Message: "",
			},
			want: "401: ",
		},
		{
			name: "Success code",
			err: Error{
				Code:    Success,
				Message: "Operation completed",
			},
			want: "200: Operation completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ShouldRetry(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want bool
	}{
		// 1xx - Informational (should not retry)
		{
			name: "Continue should not retry",
			err:  Error{Code: Continue},
			want: false,
		},
		{
			name: "SwitchingProtocols should not retry",
			err:  Error{Code: SwitchingProtocols},
			want: false,
		},
		{
			name: "Processing should not retry",
			err:  Error{Code: Processing},
			want: false,
		},
		{
			name: "EarlyHints should not retry",
			err:  Error{Code: EarlyHints},
			want: false,
		},

		// 2xx - Success (should not retry)
		{
			name: "Success should not retry",
			err:  Error{Code: Success},
			want: false,
		},
		{
			name: "Created should not retry",
			err:  Error{Code: Created},
			want: false,
		},
		{
			name: "Accepted should not retry",
			err:  Error{Code: Accepted},
			want: false,
		},

		// 4xx - Client Error (should not retry)
		{
			name: "InvalidRequest should not retry",
			err:  Error{Code: InvalidRequest},
			want: false,
		},
		{
			name: "Unauthorized should not retry",
			err:  Error{Code: Unauthorized},
			want: false,
		},
		{
			name: "Forbidden should not retry",
			err:  Error{Code: Forbidden},
			want: false,
		},
		{
			name: "NotFound should not retry",
			err:  Error{Code: NotFound},
			want: false,
		},
		{
			name: "MethodNotAllowed should not retry",
			err:  Error{Code: MethodNotAllowed},
			want: false,
		},
		{
			name: "RequestTimeout should not retry",
			err:  Error{Code: RequestTimeout},
			want: false,
		},
		{
			name: "Conflict should not retry",
			err:  Error{Code: Conflict},
			want: false,
		},
		{
			name: "TooManyRequests should not retry",
			err:  Error{Code: TooManyRequests},
			want: false,
		},

		// 5xx - Server Error (should retry)
		{
			name: "InternalServer should retry",
			err:  Error{Code: InternalServer},
			want: true,
		},
		{
			name: "NotImplemented should retry",
			err:  Error{Code: NotImplemented},
			want: true,
		},
		{
			name: "BadGateway should retry",
			err:  Error{Code: BadGateway},
			want: true,
		},
		{
			name: "ServiceUnavailable should retry",
			err:  Error{Code: ServiceUnavailable},
			want: true,
		},
		{
			name: "GatewayTimeout should retry",
			err:  Error{Code: GatewayTimeout},
			want: true,
		},
		{
			name: "HTTPVersionError should retry",
			err:  Error{Code: HTTPVersionError},
			want: true,
		},
		{
			name: "InsufficientStorage should retry",
			err:  Error{Code: InsufficientStorage},
			want: true,
		},
		{
			name: "NetworkAuthenticationRequired should retry",
			err:  Error{Code: NetworkAuthenticationRequired},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.ShouldRetry(); got != tt.want {
				t.Errorf("Error.ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorCode_Values(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		// 1xx Informational
		{"Continue", Continue, 100},
		{"SwitchingProtocols", SwitchingProtocols, 101},
		{"Processing", Processing, 102},
		{"EarlyHints", EarlyHints, 103},

		// 2xx Success
		{"Success", Success, 200},
		{"Created", Created, 201},
		{"Accepted", Accepted, 202},
		{"NonAuthoritativeInformation", NonAuthoritativeInformation, 203},
		{"NoContent", NoContent, 204},
		{"ResetContent", ResetContent, 205},
		{"PartialContent", PartialContent, 206},
		{"MultiStatus", MultiStatus, 207},
		{"AlreadyReported", AlreadyReported, 208},
		{"IMUsed", IMUsed, 226},

		// 4xx Client Error
		{"InvalidRequest", InvalidRequest, 400},
		{"Unauthorized", Unauthorized, 401},
		{"PaymentRequired", PaymentRequired, 402},
		{"Forbidden", Forbidden, 403},
		{"NotFound", NotFound, 404},
		{"MethodNotAllowed", MethodNotAllowed, 405},
		{"NotAcceptable", NotAcceptable, 406},
		{"ProxyAuthenticationRequired", ProxyAuthenticationRequired, 407},
		{"RequestTimeout", RequestTimeout, 408},
		{"Conflict", Conflict, 409},
		{"Gone", Gone, 410},
		{"LengthRequired", LengthRequired, 411},
		{"PreconditionFailed", PreconditionFailed, 412},
		{"PayloadTooLarge", PayloadTooLarge, 413},
		{"URITooLong", URITooLong, 414},
		{"UnsupportedMediaType", UnsupportedMediaType, 415},
		{"RangeNotSatisfiable", RangeNotSatisfiable, 416},
		{"ExpectationFailed", ExpectationFailed, 417},
		{"ImATeapot", ImATeapot, 418},
		{"MisdirectedRequest", MisdirectedRequest, 421},
		{"UnprocessableEntity", UnprocessableEntity, 422},
		{"Locked", Locked, 423},
		{"FailedDependency", FailedDependency, 424},
		{"TooEarly", TooEarly, 425},
		{"UpgradeRequired", UpgradeRequired, 426},
		{"PreconditionRequired", PreconditionRequired, 428},
		{"TooManyRequests", TooManyRequests, 429},
		{"RequestHeaderFieldsTooLarge", RequestHeaderFieldsTooLarge, 431},
		{"UnavailableForLegalReasons", UnavailableForLegalReasons, 451},

		// 5xx Server Error
		{"InternalServer", InternalServer, 500},
		{"NotImplemented", NotImplemented, 501},
		{"BadGateway", BadGateway, 502},
		{"ServiceUnavailable", ServiceUnavailable, 503},
		{"GatewayTimeout", GatewayTimeout, 504},
		{"HTTPVersionError", HTTPVersionError, 505},
		{"VariantAlsoNegotiates", VariantAlsoNegotiates, 506},
		{"InsufficientStorage", InsufficientStorage, 507},
		{"LoopDetected", LoopDetected, 508},
		{"NotExtended", NotExtended, 510},
		{"NetworkAuthenticationRequired", NetworkAuthenticationRequired, 511},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.code) != tt.expected {
				t.Errorf("ErrorCode %s = %v, want %v", tt.name, tt.code, tt.expected)
			}
		})
	}
}

func TestError_EdgeCases(t *testing.T) {
	t.Run("Nil details map", func(t *testing.T) {
		err := Throw(NotFound, "test", nil)
		if err.Details != nil {
			t.Errorf("Expected nil details, got %v", err.Details)
		}
	})

	t.Run("Empty details map", func(t *testing.T) {
		err := Throw(NotFound, "test", map[string]interface{}{})
		if err.Details == nil {
			t.Errorf("Expected empty map, got nil")
		}
		if len(err.Details) != 0 {
			t.Errorf("Expected empty map, got %v", err.Details)
		}
	})

	t.Run("Complex details", func(t *testing.T) {
		details := map[string]interface{}{
			"string":  "value",
			"number":  42,
			"boolean": true,
			"nested": map[string]interface{}{
				"key": "value",
			},
		}
		err := Throw(InternalServer, "complex error", details)

		if err.Details == nil {
			t.Fatal("Expected details, got nil")
		}

		// Test basic types
		if err.Details["string"] != "value" {
			t.Errorf("Details[string] = %v, want %v", err.Details["string"], "value")
		}
		if err.Details["number"] != 42 {
			t.Errorf("Details[number] = %v, want %v", err.Details["number"], 42)
		}
		if err.Details["boolean"] != true {
			t.Errorf("Details[boolean] = %v, want %v", err.Details["boolean"], true)
		}

		// Test nested map
		nested, ok := err.Details["nested"].(map[string]interface{})
		if !ok {
			t.Errorf("Details[nested] is not a map, got %T", err.Details["nested"])
		} else if nested["key"] != "value" {
			t.Errorf("Details[nested][key] = %v, want %v", nested["key"], "value")
		}
	})
}

func TestError_RetryLogic(t *testing.T) {
	// Test the retry logic for all error codes
	shouldNotRetry := []ErrorCode{
		Continue, SwitchingProtocols, Processing, EarlyHints,
		Success, Created, Accepted, NonAuthoritativeInformation,
		NoContent, ResetContent, PartialContent, MultiStatus,
		AlreadyReported, IMUsed,
		InvalidRequest, Unauthorized, PaymentRequired, Forbidden,
		NotFound, MethodNotAllowed, NotAcceptable,
		ProxyAuthenticationRequired, RequestTimeout, Conflict,
		Gone, LengthRequired, PreconditionFailed, PayloadTooLarge,
		URITooLong, UnsupportedMediaType, RangeNotSatisfiable,
		ExpectationFailed, ImATeapot, MisdirectedRequest,
		UnprocessableEntity, Locked, FailedDependency, TooEarly,
		UpgradeRequired, PreconditionRequired, TooManyRequests,
		RequestHeaderFieldsTooLarge, UnavailableForLegalReasons,
	}

	shouldRetry := []ErrorCode{
		InternalServer, NotImplemented, BadGateway, ServiceUnavailable,
		GatewayTimeout, HTTPVersionError, VariantAlsoNegotiates,
		InsufficientStorage, LoopDetected, NotExtended,
		NetworkAuthenticationRequired,
	}

	for _, code := range shouldNotRetry {
		t.Run("ShouldNotRetry_"+fmt.Sprintf("%d", code), func(t *testing.T) {
			err := Error{Code: code}
			if err.ShouldRetry() {
				t.Errorf("Error code %v should not retry, but ShouldRetry() returned true", code)
			}
		})
	}

	for _, code := range shouldRetry {
		t.Run("ShouldRetry_"+fmt.Sprintf("%d", code), func(t *testing.T) {
			err := Error{Code: code}
			if !err.ShouldRetry() {
				t.Errorf("Error code %v should retry, but ShouldRetry() returned false", code)
			}
		})
	}
}

func TestError_StringRepresentation(t *testing.T) {
	tests := []struct {
		name     string
		err      Error
		expected string
	}{
		{
			name: "Simple error",
			err: Error{
				Code:    NotFound,
				Message: "User not found",
			},
			expected: "404: User not found",
		},
		{
			name: "Error with special characters",
			err: Error{
				Code:    InvalidRequest,
				Message: "Invalid input: \"test\" & 'value'",
			},
			expected: "400: Invalid input: \"test\" & 'value'",
		},
		{
			name: "Error with unicode",
			err: Error{
				Code:    InternalServer,
				Message: "Database error: 数据库连接失败",
			},
			expected: "500: Database error: 数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
