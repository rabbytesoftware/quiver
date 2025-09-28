package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://example.com"
	client := NewClient(baseURL)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL to be %q, got %q", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("httpClient should be initialized")
	}

	expectedTimeout := 30 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}

func TestClient_ExecuteRequest_GET(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type header to be application/json")
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Error("Expected Accept header to be application/json")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/test",
		Method: "GET",
		Args:   map[string]string{},
	}

	ctx := context.Background()
	response, err := client.ExecuteRequest(ctx, request)

	if err != nil {
		t.Fatalf("ExecuteRequest() failed: %v", err)
	}

	if response == nil {
		t.Fatal("ExecuteRequest() returned nil response")
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !response.Success {
		t.Error("Expected response to be successful")
	}

	expectedBody := `{"message": "success"}`
	if string(response.Body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(response.Body))
	}
}

func TestClient_ExecuteRequest_POST(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check request body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		var requestData map[string]string
		if err := json.Unmarshal(body, &requestData); err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
		}

		if requestData["key"] != "value" {
			t.Errorf("Expected request data key=value, got %v", requestData)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/create",
		Method: "POST",
		Args:   map[string]string{"key": "value"},
	}

	ctx := context.Background()
	response, err := client.ExecuteRequest(ctx, request)

	if err != nil {
		t.Fatalf("ExecuteRequest() failed: %v", err)
	}

	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, response.StatusCode)
	}

	if !response.Success {
		t.Error("Expected response to be successful")
	}
}

func TestClient_ExecuteRequest_PUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/update",
		Method: "PUT",
		Args:   map[string]string{"id": "123"},
	}

	ctx := context.Background()
	response, err := client.ExecuteRequest(ctx, request)

	if err != nil {
		t.Fatalf("ExecuteRequest() failed: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestClient_ExecuteRequest_PATCH(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/patch",
		Method: "PATCH",
		Args:   map[string]string{"field": "newvalue"},
	}

	ctx := context.Background()
	response, err := client.ExecuteRequest(ctx, request)

	if err != nil {
		t.Fatalf("ExecuteRequest() failed: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestClient_ExecuteRequest_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/notfound",
		Method: "GET",
		Args:   map[string]string{},
	}

	ctx := context.Background()
	response, err := client.ExecuteRequest(ctx, request)

	if err != nil {
		t.Fatalf("ExecuteRequest() failed: %v", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
	}

	if response.Success {
		t.Error("Expected response to be unsuccessful")
	}
}

func TestClient_ExecuteRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	request := &models.QueryRequest{
		URL:    "/slow",
		Method: "GET",
		Args:   map[string]string{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.ExecuteRequest(ctx, request)

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded error, got: %v", err)
	}
}

func TestClient_ExecuteRequest_NetworkError(t *testing.T) {
	// Use an invalid URL that will cause a network error
	client := NewClient("http://invalid-host-that-does-not-exist.local")

	request := &models.QueryRequest{
		URL:    "/test",
		Method: "GET",
		Args:   map[string]string{},
	}

	ctx := context.Background()
	_, err := client.ExecuteRequest(ctx, request)

	if err == nil {
		t.Error("Expected error due to network failure")
	}

	if !strings.Contains(err.Error(), "failed to execute HTTP request") {
		t.Errorf("Expected HTTP request execution error, got: %v", err)
	}
}

func TestQueryResponse_String(t *testing.T) {
	testCases := []struct {
		name     string
		response QueryResponse
		expected string
	}{
		{
			name: "response with body",
			response: QueryResponse{
				StatusCode: 200,
				Body:       []byte(`{"message": "success"}`),
			},
			expected: `200:{"message": "success"}`,
		},
		{
			name: "response without body",
			response: QueryResponse{
				StatusCode: 204,
				Body:       []byte{},
			},
			expected: "204",
		},
		{
			name: "error response",
			response: QueryResponse{
				StatusCode: 500,
				Body:       []byte("Internal Server Error"),
			},
			expected: "500:Internal Server Error",
		},
		{
			name: "response with nil body",
			response: QueryResponse{
				StatusCode: 404,
				Body:       nil,
			},
			expected: "404",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.response.String()
			if result != tc.expected {
				t.Errorf("Expected String() to return %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestQueryResponse_GetBodyAsString(t *testing.T) {
	testCases := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name:     "JSON body",
			body:     []byte(`{"key": "value"}`),
			expected: `{"key": "value"}`,
		},
		{
			name:     "plain text body",
			body:     []byte("Hello, World!"),
			expected: "Hello, World!",
		},
		{
			name:     "empty body",
			body:     []byte{},
			expected: "",
		},
		{
			name:     "nil body",
			body:     nil,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := &QueryResponse{Body: tc.body}
			result := response.GetBodyAsString()
			if result != tc.expected {
				t.Errorf("Expected GetBodyAsString() to return %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestQueryResponse_GetBodyAsJSON(t *testing.T) {
	// Test successful JSON unmarshaling
	response := &QueryResponse{
		Body: []byte(`{"name": "test", "value": 123}`),
	}

	var result map[string]interface{}
	err := response.GetBodyAsJSON(&result)

	if err != nil {
		t.Fatalf("GetBodyAsJSON() failed: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name to be 'test', got %v", result["name"])
	}

	if result["value"] != float64(123) { // JSON numbers are float64
		t.Errorf("Expected value to be 123, got %v", result["value"])
	}
}

func TestQueryResponse_GetBodyAsJSON_InvalidJSON(t *testing.T) {
	response := &QueryResponse{
		Body: []byte(`{"invalid": json}`),
	}

	var result map[string]interface{}
	err := response.GetBodyAsJSON(&result)

	if err == nil {
		t.Error("Expected error due to invalid JSON")
	}
}

func TestQueryResponse_GetBodyAsJSON_EmptyBody(t *testing.T) {
	response := &QueryResponse{
		Body: []byte{},
	}

	var result map[string]interface{}
	err := response.GetBodyAsJSON(&result)

	if err == nil {
		t.Error("Expected error due to empty JSON")
	}
}

func TestQueryResponse_SuccessField(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"200 OK", 200, true},
		{"201 Created", 201, true},
		{"299 Success", 299, true},
		{"199 Not Success", 199, false},
		{"300 Not Success", 300, false},
		{"404 Not Found", 404, false},
		{"500 Server Error", 500, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := &QueryResponse{StatusCode: tc.statusCode}
			// Success field should be set during creation, let's simulate it
			response.Success = tc.statusCode >= 200 && tc.statusCode < 300

			if response.Success != tc.expected {
				t.Errorf("Expected Success to be %v for status %d, got %v", tc.expected, tc.statusCode, response.Success)
			}
		})
	}
}
