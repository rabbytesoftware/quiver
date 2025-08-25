# Gin Framework Migration Guide

## Overview

The Quiver server has been migrated from Gorilla Mux to Gin framework, providing better performance, middleware support, and developer experience. This guide covers the migration changes and new API structure.

## Migration Summary

### Framework Changes
- **Before**: Gorilla Mux + net/http
- **After**: Gin Gonic framework
- **Benefits**: Better performance, built-in middleware, cleaner routing, automatic JSON binding

### Package Structure Changes
- **Removed**: `internal/server/response/` (old Gorilla Mux style)
- **Added**: `internal/server/response/` (Gin-optimized response utilities)
- **Updated**: All handler packages to use Gin context

## Key Changes

### 1. Request Handling

#### Before (Gorilla Mux)
```go
func (h *Handler) GetArrow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    
    // Manual JSON parsing
    var req ArrowRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
}
```

#### After (Gin)
```go
func (h *Handler) GetArrow(c *gin.Context) {
    name := c.Param("name")
    
    // Automatic JSON binding
    var req ArrowRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request body", err.Error())
        return
    }
}
```

### 2. Response Handling

#### Before (Manual JSON)
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "data": arrow,
})
```

#### After (Standardized Response)
```go
response.Success(c, "Arrow retrieved successfully", arrow)
```

### 3. Route Parameter Extraction

#### Before (Gorilla Mux)
```go
vars := mux.Vars(r)
name := vars["name"]
id := vars["id"]
```

#### After (Gin)
```go
name := c.Param("name")
id := c.Param("id")
```

### 4. Query Parameter Handling

#### Before (net/http)
```go
query := r.URL.Query().Get("q")
limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
```

#### After (Gin)
```go
query := c.Query("q")
limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
```

## New API Structure

### Health Endpoints (Root Level)
```
GET /health          # Health check
GET /ready           # Readiness probe
GET /live            # Liveness probe
```

### API v1 Endpoints

#### Arrows Management
```
GET    /api/v1/arrows/search                    # Search arrows
POST   /api/v1/arrows/:name/install            # Install arrow
POST   /api/v1/arrows/:name/execute            # Execute arrow
DELETE /api/v1/arrows/:name/uninstall          # Uninstall arrow
PUT    /api/v1/arrows/:name/update             # Update arrow
POST   /api/v1/arrows/:name/validate           # Validate arrow
GET    /api/v1/arrows/installed                # List installed arrows
GET    /api/v1/arrows/:name/status             # Get arrow status
GET    /api/v1/arrows/status/:status           # Get arrows by status
GET    /api/v1/arrows/:name                    # Get arrow info
GET    /api/v1/arrows/statuses                 # List arrow statuses
```

#### Package Management
```
GET    /api/v1/packages/                       # List packages
GET    /api/v1/packages/:id                    # Get package info
POST   /api/v1/packages/:id/start              # Start package
POST   /api/v1/packages/:id/stop               # Stop package
GET    /api/v1/packages/:id/status             # Get package status
GET    /api/v1/packages/installed              # List installed packages
GET    /api/v1/packages/status/:status         # Get packages by status
GET    /api/v1/packages/statuses               # List package statuses
```

#### Repository Management
```
POST   /api/v1/repositories/                   # Add repository
DELETE /api/v1/repositories/                   # Remove repository
GET    /api/v1/repositories/                   # List repositories
GET    /api/v1/repositories/:repository        # Get repository info
GET    /api/v1/repositories/search             # Search repositories
POST   /api/v1/repositories/:repository/validate # Validate repository
```

#### Server Management
```
GET    /api/v1/server/info                     # Server information
GET    /api/v1/server/status                   # Server status
GET    /api/v1/server/health                   # Server health
GET    /api/v1/server/metrics                  # Server metrics
GET    /api/v1/server/version                  # Server version
```

#### Netbridge Management
```
POST   /api/v1/netbridge/open                  # Open port
POST   /api/v1/netbridge/close                 # Close port
POST   /api/v1/netbridge/port/:port/open       # Open port by URL
DELETE /api/v1/netbridge/port/:port            # Close port by URL
GET    /api/v1/netbridge/ports                 # List open ports
POST   /api/v1/netbridge/refresh               # Refresh public IP
GET    /api/v1/netbridge/status                # Get netbridge status
POST   /api/v1/netbridge/auto                  # Auto open port
POST   /api/v1/netbridge/auto/:port            # Auto open port by URL
```

## Request/Response Examples

### Arrow Installation

#### Request
```bash
curl -X POST http://localhost:8080/api/v1/arrows/cs2-server/install \
  -H "Content-Type: application/json" \
  -d '{
    "variables": {
      "SERVER_NAME": "My CS2 Server",
      "MAX_PLAYERS": 12
    }
  }'
```

#### Response
```json
{
  "success": true,
  "message": "Arrow installed successfully",
  "data": {
    "name": "cs2-server",
    "version": "1.0.0",
    "status": "installed",
    "installedAt": "2024-01-01T12:00:00Z"
  }
}
```

### Package Status

#### Request
```bash
curl -X GET http://localhost:8080/api/v1/packages/cs2-server/status
```

#### Response
```json
{
  "success": true,
  "message": "Package status retrieved successfully",
  "data": {
    "id": "cs2-server",
    "status": "running",
    "uptime": "2h30m",
    "lastUpdate": "2024-01-01T12:00:00Z"
  }
}
```

### Error Response

#### Request
```bash
curl -X GET http://localhost:8080/api/v1/arrows/nonexistent
```

#### Response
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Arrow not found"
  }
}
```

## Middleware Integration

### CORS Middleware
```go
func (s *Server) setupMiddleware() {
    s.gin.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

### Request Logging
```go
func (s *Server) setupMiddleware() {
    s.gin.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC3339),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    }))
}
```

## Benefits of Migration

### 1. **Performance Improvements**
- Faster request routing
- Reduced memory allocation
- Better HTTP/2 support
- Built-in JSON binding and validation

### 2. **Developer Experience**
- Cleaner, more readable code
- Automatic JSON binding
- Better error handling
- Standardized response format

### 3. **Middleware Ecosystem**
- Built-in recovery middleware
- CORS support
- Request logging
- Authentication middleware

### 4. **Testing Improvements**
- Easier unit testing with Gin test context
- Better mock support
- Cleaner test setup

## Backward Compatibility

All existing API endpoints remain functional with the same:
- URL patterns
- Request/response formats
- HTTP methods
- Authentication requirements

The migration maintains 100% backward compatibility while improving performance and maintainability.

## Development Guidelines

### Adding New Endpoints
1. Add handler method in appropriate module
2. Add route in module's `routes.go`
3. Use standardized response functions
4. Include proper error handling

### Request Validation
```go
func (h *Handler) CreateArrow(c *gin.Context) {
    var req CreateArrowRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, map[string]string{
            "body": "Invalid JSON format",
        })
        return
    }
    
    // Validate required fields
    if req.Name == "" {
        response.BadRequest(c, "Arrow name is required")
        return
    }
    
    // Process request...
}
```

### Error Handling
```go
func (h *Handler) ProcessRequest(c *gin.Context) {
    result, err := h.service.Process()
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            response.NotFound(c, "Resource")
        case errors.Is(err, ErrValidation):
            response.BadRequest(c, "Validation failed", err.Error())
        default:
            response.InternalServerError(c, "Processing failed", err.Error())
        }
        return
    }
    
    response.Success(c, "Request processed successfully", result)
}
```

This migration provides a solid foundation for future API development while maintaining compatibility and improving performance. 