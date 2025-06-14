# Server Package

This package contains the HTTP server implementation for the Quiver application, organized following Go best practices with clear separation of concerns.

## Package Structure

```
internal/server/
├── server.go      # Main server struct and lifecycle management
├── handlers.go    # HTTP request handlers
├── routes.go      # Route definitions and setup
├── middleware.go  # HTTP middleware functions
├── response.go    # Response utilities and helpers
└── README.md      # This documentation
```

## Architecture Overview

### Core Components

1. **Server (`server.go`)**
   - Main server struct and lifecycle management
   - Server configuration and startup/shutdown logic
   - Server factory functions

2. **Handlers (`handlers.go`)**
   - All HTTP request handlers
   - Business logic for each endpoint
   - Request validation and response formatting

3. **Routes (`routes.go`)**
   - Route definitions and URL patterns
   - Route grouping and organization
   - Middleware attachment to specific routes

4. **Middleware (`middleware.go`)**
   - Cross-cutting concerns (logging, CORS, recovery)
   - Request/response interceptors
   - Security and authentication middleware

5. **Response Utilities (`response.go`)**
   - Common response helpers
   - Standardized error formatting
   - JSON response utilities

## Usage

### Creating a Server

```go
import (
    "github.com/rabbytesoftware/quiver/internal/config"
    "github.com/rabbytesoftware/quiver/internal/logger"
    "github.com/rabbytesoftware/quiver/internal/server"
    "github.com/rabbytesoftware/quiver/packages"
)

// Create server configuration
cfg := config.ServerConfig{
    Host: "0.0.0.0",
    Port: 8080,
    ReadTimeout: 30,
    WriteTimeout: 30,
}

// Initialize dependencies
logger := logger.New(loggerConfig)
pkgManager := packages.NewArrowsServer("./pkgs")

// Create server
srv := server.New(cfg, pkgManager, logger)

// Start server
ctx := context.Background()
err := srv.Start(ctx)
```

### Adding New Endpoints

1. **Add handler function in `handlers.go`:**
```go
func (h *Handlers) NewEndpointHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation here
    h.writeJSON(w, http.StatusOK, data)
}
```

2. **Add route in `routes.go`:**
```go
func (s *Server) setupRoutes() {
    // ... existing routes
    api.HandleFunc("/new-endpoint", s.handlers.NewEndpointHandler).Methods("GET")
}
```

### Adding Middleware

1. **Add middleware function in `middleware.go`:**
```go
func (s *Server) newMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware logic here
        next.ServeHTTP(w, r)
    })
}
```

2. **Register middleware in `server.go`:**
```go
func (s *Server) setupMiddleware() {
    // ... existing middleware
    s.router.Use(s.newMiddleware)
}
```

## API Endpoints

The server provides the following REST API endpoints:

### Health Check
- `GET /health` - Server health status

### Package Management
- `GET /api/v1/packages` - List all packages
- `GET /api/v1/packages/{id}` - Get specific package
- `POST /api/v1/packages/{id}/start` - Start a package
- `POST /api/v1/packages/{id}/stop` - Stop a package
- `GET /api/v1/packages/{id}/status` - Get package status

### Server Management
- `GET /api/v1/server/info` - Server information
- `GET /api/v1/server/status` - Server status

## Design Principles

### Separation of Concerns
- **Server**: Handles server lifecycle and coordination
- **Handlers**: Contains business logic for each endpoint
- **Routes**: Defines URL patterns and routing logic
- **Middleware**: Handles cross-cutting concerns
- **Response**: Provides consistent response formatting

### Single Responsibility
Each file has a specific purpose and responsibility, making the codebase easier to:
- Understand and navigate
- Test individual components
- Modify without affecting other components
- Maintain and debug

### Testability
The modular structure makes it easy to:
- Unit test individual handlers
- Mock dependencies
- Test middleware in isolation
- Integration test the entire server

### Extensibility
Adding new functionality is straightforward:
- New endpoints: Add handler + route
- New middleware: Add function + register
- New response types: Extend response utilities

## Best Practices

1. **Error Handling**: Use the response utilities for consistent error formatting
2. **Logging**: Log important events and errors using the structured logger
3. **Validation**: Validate input data in handlers before processing
4. **Security**: Use appropriate middleware for authentication and authorization
5. **Documentation**: Keep this README updated when adding new functionality

## Testing

To test the server components:

```bash
# Run all server tests
go test ./internal/server/...

# Run with coverage
go test -cover ./internal/server/...

# Run specific test file
go test ./internal/server/handlers_test.go
```

## Dependencies

- **gorilla/mux**: HTTP router and URL matcher
- **internal/config**: Configuration management
- **internal/logger**: Structured logging
- **internal/ui**: User interface components
- **packages**: Package management system

## Configuration

The server is configured through the `config.ServerConfig` struct:

```go
type ServerConfig struct {
    Port         int    // Server port (default: 8080)
    Host         string // Server host (default: "0.0.0.0")
    ReadTimeout  int    // Read timeout in seconds
    WriteTimeout int    // Write timeout in seconds
}
```

Environment variables:
- `QUIVER_PORT`: Override server port
- `QUIVER_HOST`: Override server host
- `QUIVER_LOG_LEVEL`: Set logging level

## Security Considerations

- CORS middleware is enabled for cross-origin requests
- Recovery middleware prevents panics from crashing the server
- Input validation should be performed in handlers
- Consider adding authentication middleware for protected endpoints

## Performance

- Graceful shutdown ensures proper cleanup
- Configurable timeouts prevent hanging connections
- Structured logging minimizes performance impact
- Middleware pipeline is optimized for minimal overhead 