package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// JSON writes a standardized JSON response
func JSON(c *gin.Context, statusCode int, data interface{}) {
	response := StandardResponse{
		Success:   statusCode >= 200 && statusCode < 300,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(statusCode, response)
}

// Success writes a success response with data
func Success(c *gin.Context, message string, data interface{}) {
	response := StandardResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// Created writes a 201 created response
func Created(c *gin.Context, message string, data interface{}) {
	response := StandardResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, response)
}

// Error writes an error response with detailed information
func Error(c *gin.Context, statusCode int, code, message string, details ...string) {
	errorInfo := &ErrorInfo{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		errorInfo.Details = details[0]
	}

	response := StandardResponse{
		Success:   false,
		Error:     errorInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(statusCode, response)
}

// BadRequest writes a 400 bad request response
func BadRequest(c *gin.Context, message string, details ...string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// NotFound writes a 404 not found response
func NotFound(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", resource+" not found")
}

// InternalServerError writes a 500 internal server error response
func InternalServerError(c *gin.Context, message string, details ...string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details...)
}

// Unauthorized writes a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden writes a 403 forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

// Conflict writes a 409 conflict response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message)
}

// HealthCheck writes a health check response
func HealthCheck(c *gin.Context, service, version string) {
	healthData := gin.H{
		"status":  "ok",
		"service": service,
		"version": version,
	}

	Success(c, "Service is healthy", healthData)
}

// Pagination contains pagination metadata
type Pagination struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Total    int `json:"total"`
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
}

// PaginatedResponse writes a paginated response
func PaginatedResponse(c *gin.Context, data interface{}, pagination Pagination) {
	responseData := gin.H{
		"items":      data,
		"pagination": pagination,
	}

	Success(c, "Data retrieved successfully", responseData)
}

// ValidationError writes a validation error response
func ValidationError(c *gin.Context, errors map[string]string) {
	response := StandardResponse{
		Success:   false,
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
		},
		Data:      gin.H{"validation_errors": errors},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusUnprocessableEntity, response)
} 