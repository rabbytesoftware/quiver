# Standardized Response System

## Overview

The Quiver server uses a standardized response system built on top of Gin framework. This system provides consistent API responses, error handling, and JSON formatting across all endpoints.

## Architecture

### Response Package Structure

```
internal/server/response/
└── response.go     # Standardized response utilities
```

### Core Components

#### 1. StandardResponse Structure
```go
type StandardResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}
```

#### 2. ErrorInfo Structure
```go
type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

## Response Functions

### Success Responses

#### `Success(c *gin.Context, message string, data interface{})`
Returns a successful response with data.

```go
response.Success(c, "Arrow installed successfully", map[string]interface{}{
    "name": "cs2-server",
    "version": "1.0.0",
    "status": "installed",
})
```

**Response:**
```json
{
    "success": true,
    "message": "Arrow installed successfully",
    "data": {
        "name": "cs2-server",
        "version": "1.0.0",
        "status": "installed"
    }
}
```

#### `Created(c *gin.Context, message string, data interface{})`
Returns a 201 Created response for resource creation.

```go
response.Created(c, "Repository added successfully", map[string]string{
    "repository": "https://github.com/example/arrows",
    "status": "added",
})
```

### Error Responses

#### `Error(c *gin.Context, statusCode int, message, details string)`
Returns a generic error response.

```go
response.Error(c, http.StatusInternalServerError, "Database connection failed", "Connection timeout after 30s")
```

**Response:**
```json
{
    "success": false,
    "error": {
        "code": "INTERNAL_SERVER_ERROR",
        "message": "Database connection failed",
        "details": "Connection timeout after 30s"
    }
}
```

#### `BadRequest(c *gin.Context, message string, details ...string)`
Returns a 400 Bad Request response.

```go
response.BadRequest(c, "Invalid arrow name", "Arrow name must be alphanumeric")
```

#### `NotFound(c *gin.Context, resource string)`
Returns a 404 Not Found response.

```go
response.NotFound(c, "Arrow")
```

**Response:**
```json
{
    "success": false,
    "error": {
        "code": "NOT_FOUND",
        "message": "Arrow not found"
    }
}
```

#### `InternalServerError(c *gin.Context, message string, details ...string)`
Returns a 500 Internal Server Error response.

```go
response.InternalServerError(c, "Failed to process request", "Database query failed")
```

### Specialized Responses

#### `HealthCheck(c *gin.Context, service, version string)`
Returns a standardized health check response.

```go
response.HealthCheck(c, "Quiver", "1.0.0")
```

**Response:**
```json
{
    "success": true,
    "message": "Service is healthy",
    "data": {
        "service": "Quiver",
        "version": "1.0.0",
        "status": "healthy",
        "timestamp": "2024-01-01T12:00:00Z"
    }
}
```

#### `PaginatedResponse(c *gin.Context, message string, data interface{}, pagination PaginationInfo)`
Returns a paginated response with metadata.

```go
response.PaginatedResponse(c, "Arrows retrieved successfully", arrows, response.PaginationInfo{
    Page:       1,
    PerPage:    10,
    Total:      50,
    TotalPages: 5,
})
```

#### `ValidationError(c *gin.Context, errors map[string]string)`
Returns a 422 validation error response.

```go
response.ValidationError(c, map[string]string{
    "name": "Name is required",
    "version": "Version must be a valid semver",
})
```

## Usage Examples

### In Arrow Handler

```go
func (h *Handler) InstallArrow(c *gin.Context) {
    name := c.Param("name")
    if name == "" {
        response.BadRequest(c, "Arrow name is required")
        return
    }

    arrow, err := h.arrowManager.Install(name)
    if err != nil {
        response.InternalServerError(c, "Failed to install arrow", err.Error())
        return
    }

    response.Created(c, "Arrow installed successfully", map[string]interface{}{
        "name":        arrow.Name,
        "version":     arrow.Version,
        "status":      arrow.Status,
        "installedAt": arrow.InstalledAt,
    })
}
```

### In Package Handler

```go
func (h *Handler) GetPackageStatus(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.BadRequest(c, "Package ID is required")
        return
    }

    status, err := h.packageManager.GetStatus(id)
    if err != nil {
        response.NotFound(c, "Package")
        return
    }

    response.Success(c, "Package status retrieved successfully", map[string]interface{}{
        "id":          id,
        "status":      status.State,
        "uptime":      status.Uptime,
        "lastUpdate": status.LastUpdate,
    })
}
```

## Migration from Gorilla Mux

### Before (Gorilla Mux + Manual JSON)
```go
func (h *Handler) GetArrow(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    vars := mux.Vars(r)
    name := vars["name"]
    
    if name == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Arrow name is required",
        })
        return
    }
    
    arrow, err := h.arrowManager.Get(name)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Arrow not found",
        })
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": arrow,
    })
}
```

### After (Gin + Standardized Response)
```go
func (h *Handler) GetArrow(c *gin.Context) {
    name := c.Param("name")
    if name == "" {
        response.BadRequest(c, "Arrow name is required")
        return
    }

    arrow, err := h.arrowManager.Get(name)
    if err != nil {
        response.NotFound(c, "Arrow")
        return
    }

    response.Success(c, "Arrow retrieved successfully", arrow)
}
```

## Best Practices

### 1. **Consistent Error Handling**
- Always use the standardized error response functions
- Provide meaningful error messages
- Include helpful details when appropriate

### 2. **Descriptive Success Messages**
- Use clear, action-oriented messages
- Be consistent with message patterns across similar endpoints
- Include the resource name in the message

### 3. **Structured Data Responses**
- Return consistent data structures
- Use meaningful field names
- Include relevant metadata (timestamps, IDs, etc.)

### 4. **HTTP Status Codes**
- Use appropriate HTTP status codes
- Let the response functions handle status codes
- Don't override unless necessary

### 5. **Error Details**
- Include helpful error details for debugging
- Don't expose sensitive information
- Use error codes for programmatic handling

## Error Code Standards

The response system uses standardized error codes:

- `BAD_REQUEST` - 400 Bad Request
- `NOT_FOUND` - 404 Not Found
- `INTERNAL_SERVER_ERROR` - 500 Internal Server Error
- `VALIDATION_ERROR` - 422 Unprocessable Entity
- `UNAUTHORIZED` - 401 Unauthorized
- `FORBIDDEN` - 403 Forbidden

## Response Headers

All responses include standard headers:
- `Content-Type: application/json`
- `X-Request-ID: <request-id>` (if request ID middleware is enabled)
- `X-Response-Time: <duration>` (if timing middleware is enabled)

## Testing Response Functions

```go
func TestArrowHandler_GetArrow(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    c.Params = []gin.Param{{Key: "name", Value: "test-arrow"}}
    
    handler := &Handler{arrowManager: mockArrowManager}
    handler.GetArrow(c)
    
    // Assertions
    assert.Equal(t, 200, w.Code)
    
    var response StandardResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success)
    assert.Equal(t, "Arrow retrieved successfully", response.Message)
}
```

This standardized response system ensures consistent API behavior, improved error handling, and better developer experience across all Quiver endpoints. 