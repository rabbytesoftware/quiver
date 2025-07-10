# Modular Routing Architecture

## Overview

The Quiver server has been refactored to use a modular routing architecture where each handler module manages its own routes. This approach improves maintainability, organization, and follows Go best practices for scalable web applications.

## Architecture Changes

### Before: Centralized Routing
```go
// internal/server/routes.go (100+ lines)
func (s *Server) setupRoutes() {
    // All routes defined in one place
    s.gin.GET("/health", s.handlers.Health.HealthCheck)
    api := s.gin.Group("/api/v1")
    
    // Individual setup functions for each module
    s.setupPackageRoutes(api)
    s.setupArrowRoutes(api)
    s.setupRepositoryRoutes(api)
    // ... more setup functions
}
```

### After: Modular Routing
```go
// internal/server/routes.go (15 lines)
func (s *Server) setupRoutes() {
    // Health routes (root level)
    s.handlers.Health.SetupRoutes(s.gin.Group(""))
    
    // API v1 routes
    api := s.gin.Group("/api/v1")
    
    // Each module manages its own routes
    s.handlers.Arrows.SetupRoutes(api)
    s.handlers.Packages.SetupRoutes(api)
    s.handlers.Repositories.SetupRoutes(api)
    s.handlers.Server.SetupRoutes(api)
    s.handlers.Netbridge.SetupRoutes(api)
}
```

## Handler Module Structure

Each handler module now contains:

```
internal/server/handlers/[module]/
├── [module].go     # Handler methods and business logic
├── routes.go       # Route definitions and setup
└── types.go        # Request/response types (if needed)
```

### Example: Arrow Module

```go
// internal/server/handlers/arrows/routes.go
package arrows

import "github.com/gin-gonic/gin"

func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
    arrows := router.Group("/arrows")
    {
        // Arrow management
        arrows.GET("/search", h.SearchArrows)
        arrows.POST("/:name/install", h.InstallArrow)
        arrows.POST("/:name/execute", h.ExecuteArrow)
        arrows.DELETE("/:name/uninstall", h.UninstallArrow)
        arrows.PUT("/:name/update", h.UpdateArrow)
        arrows.POST("/:name/validate", h.ValidateArrow)
        
        // Arrow information
        arrows.GET("/installed", h.GetInstalledArrows)
        arrows.GET("/:name/status", h.GetArrowStatus)
        arrows.GET("/:name", h.GetArrowInfo)
        arrows.GET("/statuses", h.ListArrowStatuses)
    }
}
```

## Benefits

### 1. **Improved Organization**
- Each module is self-contained
- Routes are co-located with their handlers
- Easier to understand module boundaries

### 2. **Better Maintainability**
- Adding new endpoints only requires editing the relevant module
- No need to modify central routing configuration
- Reduced merge conflicts

### 3. **Enhanced Testability**
- Routes can be tested in isolation
- Module-specific route testing
- Clear separation of concerns

### 4. **Scalability**
- Easy to add new modules
- Consistent pattern for all handlers
- Better code organization as project grows

## Module Responsibilities

### Health Module (`/internal/server/handlers/health/`)
- **Routes**: `/health`, `/ready`, `/live`
- **Purpose**: Health checks and readiness probes
- **Scope**: Root-level health endpoints

### Arrows Module (`/internal/server/handlers/arrows/`)
- **Routes**: `/api/v1/arrows/*`
- **Purpose**: Arrow package management
- **Scope**: Search, install, execute, uninstall, update, validate arrows

### Packages Module (`/internal/server/handlers/packages/`)
- **Routes**: `/api/v1/packages/*`
- **Purpose**: Package lifecycle management
- **Scope**: Start, stop, status, listing packages

### Repositories Module (`/internal/server/handlers/repositories/`)
- **Routes**: `/api/v1/repositories/*`
- **Purpose**: Repository management
- **Scope**: Add, remove, search, validate repositories

### Server Module (`/internal/server/handlers/server/`)
- **Routes**: `/api/v1/server/*`
- **Purpose**: Server monitoring and information
- **Scope**: Server status, metrics, version, health

### Netbridge Module (`/internal/server/handlers/netbridge/`)
- **Routes**: `/api/v1/netbridge/*`
- **Purpose**: Network bridge management
- **Scope**: Port management, public IP, network status

## Adding New Modules

### Step 1: Create Module Directory
```bash
mkdir -p internal/server/handlers/newmodule
```

### Step 2: Create Handler
```go
// internal/server/handlers/newmodule/newmodule.go
package newmodule

import (
    "github.com/gin-gonic/gin"
    "github.com/rabbytesoftware/quiver/internal/server/response"
)

type Handler struct {
    // dependencies
}

func NewHandler(/* dependencies */) *Handler {
    return &Handler{
        // initialize dependencies
    }
}

func (h *Handler) ExampleEndpoint(c *gin.Context) {
    response.Success(c, "Example endpoint", map[string]string{
        "message": "Hello from new module",
    })
}
```

### Step 3: Create Routes
```go
// internal/server/handlers/newmodule/routes.go
package newmodule

import "github.com/gin-gonic/gin"

func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
    module := router.Group("/newmodule")
    {
        module.GET("/example", h.ExampleEndpoint)
        // Add more routes as needed
    }
}
```

### Step 4: Register in Main Server
```go
// internal/server/handlers/handlers.go
type Handlers struct {
    // ... existing handlers
    NewModule *newmodule.Handler
}

// internal/server/routes.go
func (s *Server) setupRoutes() {
    // ... existing setup
    s.handlers.NewModule.SetupRoutes(api)
}
```

## Best Practices

### 1. **Route Organization**
- Group related routes using `router.Group()`
- Use descriptive route names
- Follow RESTful conventions

### 2. **Error Handling**
- Use standardized response utilities
- Handle errors consistently across modules
- Provide meaningful error messages

### 3. **Request Validation**
- Validate input parameters
- Use Gin's binding features
- Return appropriate error responses

### 4. **Documentation**
- Document each endpoint's purpose
- Include request/response examples
- Keep route documentation up to date

## Migration Guide

### For Existing Endpoints
1. **Locate** the endpoint in the old centralized routes
2. **Move** the route definition to the appropriate module's `routes.go`
3. **Update** any route-specific logic if needed
4. **Test** the endpoint to ensure it works correctly

### For New Endpoints
1. **Identify** the appropriate module
2. **Add** the handler method to the module's handler file
3. **Add** the route to the module's `routes.go`
4. **No changes** needed to central routing configuration

This modular approach makes the codebase more maintainable, testable, and scalable while following Go best practices for web application architecture. 