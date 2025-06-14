package server

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter provides utilities for writing HTTP responses
type ResponseWriter struct{}

// NewResponseWriter creates a new response writer
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// JSON writes a JSON response with the given status code and data
func (rw *ResponseWriter) JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Error writes an error response with the given status code and message
func (rw *ResponseWriter) Error(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := map[string]string{
		"error":   message,
		"status":  http.StatusText(statusCode),
		"code":    string(rune(statusCode)),
	}
	rw.JSON(w, statusCode, errorResponse)
}

// Success writes a success response with the given message
func (rw *ResponseWriter) Success(w http.ResponseWriter, message string) {
	successResponse := map[string]string{
		"message": message,
		"status":  "success",
	}
	rw.JSON(w, http.StatusOK, successResponse)
}

// HealthCheck writes a health check response
func (rw *ResponseWriter) HealthCheck(w http.ResponseWriter, service string, version string) {
	healthResponse := map[string]interface{}{
		"status":    "ok",
		"service":   service,
		"version":   version,
		"timestamp": "2025-01-25T00:00:00Z", // You'd want to use actual timestamp
	}
	rw.JSON(w, http.StatusOK, healthResponse)
}

// NotFound writes a 404 not found response
func (rw *ResponseWriter) NotFound(w http.ResponseWriter, resource string) {
	rw.Error(w, http.StatusNotFound, resource+" not found")
}

// BadRequest writes a 400 bad request response
func (rw *ResponseWriter) BadRequest(w http.ResponseWriter, message string) {
	rw.Error(w, http.StatusBadRequest, message)
}

// InternalServerError writes a 500 internal server error response
func (rw *ResponseWriter) InternalServerError(w http.ResponseWriter, message string) {
	rw.Error(w, http.StatusInternalServerError, message)
}

// Unauthorized writes a 401 unauthorized response
func (rw *ResponseWriter) Unauthorized(w http.ResponseWriter, message string) {
	rw.Error(w, http.StatusUnauthorized, message)
}

// Forbidden writes a 403 forbidden response
func (rw *ResponseWriter) Forbidden(w http.ResponseWriter, message string) {
	rw.Error(w, http.StatusForbidden, message)
} 