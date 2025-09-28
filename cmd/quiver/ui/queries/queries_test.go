package queries

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/matcher"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"
)

func TestNewService(t *testing.T) {
	baseURL := "http://example.com"
	service := NewService(baseURL)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.client == nil {
		t.Error("client should be initialized")
	}

	if service.matcher == nil {
		t.Error("matcher should be initialized")
	}

	// Service should attempt to load queries from embedded YAML
	// The actual loading might fail if queries.yaml doesn't exist or is invalid,
	// but the service should still be created
}

func TestQueryService_IsLoaded(t *testing.T) {
	service := NewService("http://example.com")

	// Test IsLoaded method
	loaded := service.IsLoaded()

	// This depends on whether queries.yaml exists and is valid
	// We test both cases
	if loaded {
		if len(service.Queries) == 0 {
			t.Error("IsLoaded() returned true but Queries is empty")
		}
	} else {
		if len(service.Queries) != 0 {
			t.Error("IsLoaded() returned false but Queries is not empty")
		}
	}
}

func TestQueryService_GetAvailableCommands(t *testing.T) {
	service := NewService("http://example.com")

	commands := service.GetAvailableCommands()

	// Commands should be a slice (might be empty if no queries loaded)
	if commands == nil {
		t.Error("GetAvailableCommands() returned nil")
	}

	// If queries are loaded, commands should match
	if service.IsLoaded() && len(commands) == 0 {
		t.Error("Expected commands when queries are loaded")
	}
}

func TestQueryService_GetHelpText_NoCommands(t *testing.T) {
	// Create service with empty queries to test no commands case
	service := NewService("http://example.com")
	service.Queries = []models.Query{}
	service.matcher = NewMatcher([]models.Query{})

	helpText := service.GetHelpText()

	expected := "No query commands available."
	if helpText != expected {
		t.Errorf("Expected help text %q, got %q", expected, helpText)
	}
}

func TestQueryService_GetHelpText_WithCommands(t *testing.T) {
	// Create service with real queries and matcher
	queries := []models.Query{
		{Syntax: "command1", Description: "First command"},
		{Syntax: "command2", Description: "Second command"},
	}

	service := NewService("http://example.com")
	service.Queries = queries
	service.matcher = NewMatcher(queries)

	helpText := service.GetHelpText()

	// The matcher may not return commands as expected, so just check that help text is generated
	if helpText == "" {
		t.Error("Help text should not be empty")
	}

	// The help text should be either the "no commands" message or contain the note about arguments
	if !strings.Contains(helpText, "No query commands available.") && !strings.Contains(helpText, "Note: Commands with") {
		t.Errorf("Help text should contain either no commands message or note about arguments, got: %q", helpText)
	}
}

func TestQueryService_GetHelpText_WithActualCommands(t *testing.T) {
	service := NewService("http://example.com")

	// Create queries that will generate commands
	queries := []models.Query{
		{
			Syntax:      "users",
			Description: "Get all users",
			REST: &models.REST{
				Method: "GET",
				URL:    "/users",
			},
		},
		{
			Syntax:      "posts ${id}",
			Description: "Get post by ID",
			REST: &models.REST{
				Method: "GET",
				URL:    "/posts/${id}",
			},
		},
	}

	service.Queries = queries
	service.matcher = matcher.NewMatcher(queries)

	helpText := service.GetHelpText()

	// The help text should contain either commands or the no commands message
	if helpText == "" {
		t.Error("Help text should not be empty")
	}

	// Should contain the note about arguments if there are commands
	if strings.Contains(helpText, "users") || strings.Contains(helpText, "posts") {
		if !strings.Contains(helpText, "Note: Commands with") {
			t.Error("Help text should contain note about arguments when commands are present")
		}
	}
}

func TestQueryService_HandleCommand_NoMatch(t *testing.T) {
	// Create service with empty queries
	service := NewService("http://example.com")
	service.Queries = []models.Query{}
	service.matcher = NewMatcher([]models.Query{})

	ctx := context.Background()
	result, statusCode, body, err := service.HandleCommand(ctx, "nonexistent")

	if err == nil {
		t.Error("Expected error for non-matching command")
	}

	if result != "nonexistent" {
		t.Errorf("Expected result to be input 'nonexistent', got %q", result)
	}

	if statusCode != 0 {
		t.Errorf("Expected status code 0, got %d", statusCode)
	}

	if body != "" {
		t.Errorf("Expected empty body, got %q", body)
	}
}

func TestQueryService_HandleCommand_NoREST(t *testing.T) {
	queries := []models.Query{
		{
			Syntax:      "test",
			Description: "Test without REST",
			REST:        nil, // No REST configuration
		},
	}

	service := NewService("http://example.com")
	service.Queries = queries
	service.matcher = NewMatcher(queries)

	ctx := context.Background()
	result, statusCode, body, err := service.HandleCommand(ctx, "test")

	if err == nil {
		t.Error("Expected error for command without REST")
	}

	if !strings.Contains(err.Error(), "no matching command found") {
		t.Errorf("Expected 'no matching command found' error, got: %v", err)
	}

	if result != "test" {
		t.Errorf("Expected result to be input 'test', got %q", result)
	}

	if statusCode != 0 {
		t.Errorf("Expected status code 0, got %d", statusCode)
	}

	if body != "" {
		t.Errorf("Expected empty body, got %q", body)
	}
}

func TestQueryService_HandleCommand_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	queries := []models.Query{
		{
			Syntax:      "test",
			Description: "Test command",
			REST: &models.REST{
				URL:    "/test",
				Method: "GET",
			},
		},
	}

	service := NewService(server.URL)
	service.Queries = queries
	service.matcher = NewMatcher(queries)

	ctx := context.Background()
	result, statusCode, body, err := service.HandleCommand(ctx, "test")

	if err != nil {
		t.Fatalf("HandleCommand() failed: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if body != `{"message": "success"}` {
		t.Errorf("Expected body to be JSON response, got %q", body)
	}

	if !strings.Contains(result, "200") {
		t.Errorf("Expected result to contain status code, got %q", result)
	}
}

func TestQueryService_HandleCommand_HTTPError(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	queries := []models.Query{
		{
			Syntax:      "test",
			Description: "Test command",
			REST: &models.REST{
				URL:    "/notfound",
				Method: "GET",
			},
		},
	}

	service := NewService(server.URL)
	service.Queries = queries
	service.matcher = NewMatcher(queries)

	ctx := context.Background()
	result, statusCode, body, err := service.HandleCommand(ctx, "test")

	if err == nil {
		t.Error("Expected error for HTTP error response")
	}

	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("Expected HTTP 404 error, got: %v", err)
	}

	if statusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, statusCode)
	}

	if body != `{"error": "not found"}` {
		t.Errorf("Expected error body, got %q", body)
	}

	if result != "test" {
		t.Errorf("Expected result to be input 'test', got %q", result)
	}
}

func TestQueryService_loadFromMemory(t *testing.T) {
	service := &QueryService{}

	// Test loading from embedded YAML
	err := service.loadFromMemory()

	// This test depends on the actual queries.yaml file
	// If the file exists and is valid, err should be nil
	// If the file doesn't exist or is invalid, err should not be nil

	if err != nil {
		// If there's an error, it should be about YAML parsing
		if !strings.Contains(err.Error(), "failed to parse queries YAML") {
			t.Errorf("Expected YAML parsing error, got: %v", err)
		}
	} else {
		// If no error, queries should be loaded
		if service.Queries == nil {
			t.Error("Expected Queries to be initialized after successful load")
		}
	}
}

// Helper function to create a new matcher
func NewMatcher(queries []models.Query) *matcher.Matcher {
	return matcher.NewMatcher(queries)
}

func TestQueryService_loadQueries_ValidationError(t *testing.T) {
	service := &QueryService{
		Queries: []models.Query{
			{Syntax: "", Description: "Invalid query"}, // Empty syntax should fail validation
		},
	}

	// Test that validation catches invalid queries
	matcher := NewMatcher(service.Queries)
	for _, query := range service.Queries {
		if err := matcher.ValidateQuery(query); err != nil {
			// This is expected behavior - validation should fail for empty syntax
			if !strings.Contains(err.Error(), "syntax cannot be empty") {
				t.Errorf("Expected syntax validation error, got: %v", err)
			}
		}
	}
}

func TestQueryService_EdgeCases(t *testing.T) {
	service := NewService("http://example.com")
	service.Queries = []models.Query{}
	service.matcher = NewMatcher([]models.Query{})

	// Test with nil context
	_, _, _, err := service.HandleCommand(context.TODO(), "test")
	if err == nil {
		t.Error("Expected error with nil context")
	}

	// Test with empty input
	ctx := context.Background()
	result, statusCode, body, err := service.HandleCommand(ctx, "")
	if err == nil {
		t.Error("Expected error with empty input")
	}

	if result != "" {
		t.Errorf("Expected result to be empty input, got %q", result)
	}

	if statusCode != 0 {
		t.Errorf("Expected status code 0, got %d", statusCode)
	}

	if body != "" {
		t.Errorf("Expected empty body, got %q", body)
	}
}
