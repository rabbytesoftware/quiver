# Project Structure

This document explains the organization and structure of the Quiver package manager codebase, following Go project layout standards and Clean Architecture principles.

## Directory Overview

```
quiver/
├── 📁 arrow.dev/              # Arrow package examples and templates
├── 📁 bin/                    # Compiled binaries (generated)
├── 📁 cmd/                    # Main applications
├── 📁 docs/                    # Documentation (this wiki)
├── 📁 internal/               # Private application code
├── 📁 logs/                    # Application logs (generated)
├── 📄 Makefile                # Build and development commands
├── 📄 go.mod                  # Go module dependencies
├── 📄 go.sum                  # Go module checksums
├── 📄 README.md               # Project overview
└── 📄 LICENSE                  # GPL-3.0 license
```

## Core Directories

### `/cmd/` - Main Applications

Contains the main entry points for the application.

```
cmd/
└── quiver/                    # Main Quiver application
    ├── main.go                # Application entry point
    ├── main_test.go           # Main function tests
    └── ui/                    # Terminal UI components
        ├── model.go           # TUI model (Bubble Tea)
        ├── model_test.go      # TUI model tests
        ├── service.go         # TUI service layer
        ├── service_test.go    # TUI service tests
        ├── domain/            # TUI domain logic
        │   ├── commands/      # Command parsing and types
        │   ├── events/        # Event system
        │   └── handlers/      # Command handlers
        ├── queries/           # API query services
        ├── services/          # TUI services
        └── styles/            # UI styling and themes
```

**Purpose**: Entry points and user interface components.

### `/internal/` - Private Application Code

Contains all private application code that should not be imported by other projects.

```
internal/
├── api/                       # REST API layer
│   ├── api.go                # API setup and configuration
│   ├── api_test.go           # API tests
│   ├── middleware/            # HTTP middleware
│   └── v1/                   # API version 1
│       ├── routes.go         # Route definitions
│       ├── routes_test.go    # Route tests
│       └── controllers/       # HTTP controllers
│           ├── arrows/       # Arrow management endpoints
│           ├── health/       # Health check endpoints
│           ├── quivers/      # Quiver management endpoints
│           └── system/       # System endpoints
├── core/                     # Core services
│   ├── core.go               # Core service initialization
│   ├── core_test.go          # Core service tests
│   ├── config/               # Configuration management
│   │   ├── config.go         # Configuration logic
│   │   ├── config_test.go    # Configuration tests
│   │   └── default.yaml      # Default configuration
│   ├── errors/               # Error handling
│   ├── metadata/             # Application metadata
│   └── watcher/              # Logging and monitoring
│       ├── service.go        # Watcher service
│       ├── service_test.go   # Watcher tests
│       ├── messages.go       # Log message types
│       └── pool/             # Connection pooling
├── infrastructure/           # Infrastructure layer
│   ├── infrastructure.go     # Infrastructure setup
│   ├── infrastructure_test.go # Infrastructure tests
│   ├── database/             # Database implementations
│   ├── fetchnshare/          # Package fetching
│   ├── netbridge/            # Network bridging
│   ├── requirements/         # System requirements
│   ├── runtime/              # Runtime management
│   ├── translator/            # Package translation
│   └── wizard/               # Setup wizard
├── models/                   # Domain models
│   ├── arrow/                # Arrow (package) model
│   │   ├── arrow.go          # Arrow entity
│   │   ├── arrow_test.go      # Arrow tests
│   │   ├── arrow-namespace.go # Arrow namespace
│   │   └── arrow-namespace_test.go
│   ├── port/                 # Port management
│   │   ├── port.go           # Port entity
│   │   ├── port_test.go      # Port tests
│   │   ├── protocol.go       # Protocol types
│   │   └── forwarding_status.go # Port forwarding status
│   ├── quiver/               # Quiver (server) model
│   │   ├── quiver.go         # Quiver entity
│   │   └── quiver_test.go    # Quiver tests
│   ├── requirement/          # System requirements
│   ├── runtime/              # Runtime methods
│   ├── system/               # System information
│   │   ├── os.go             # Operating system
│   │   ├── security.go        # Security levels
│   │   └── url.go            # URL handling
│   └── variable/             # Configuration variables
│       ├── variable.go       # Variable entity
│       ├── variable_test.go  # Variable tests
│       └── type.go           # Variable types
├── repositories/             # Data access layer
│   ├── repositories.go        # Repository container
│   ├── repositories_test.go   # Repository tests
│   ├── arrows/               # Arrow repository
│   │   ├── arrows.go         # Arrow data access
│   │   ├── arrows_test.go    # Arrow repository tests
│   │   └── interface.go      # Arrow repository interface
│   ├── quivers/              # Quiver repository
│   │   ├── quivers.go        # Quiver data access
│   │   ├── quivers_test.go   # Quiver repository tests
│   │   └── interface.go      # Quiver repository interface
│   ├── system/               # System repository
│   │   ├── system.go         # System data access
│   │   ├── system_test.go    # System repository tests
│   │   └── interface.go      # System repository interface
│   └── common/               # Common repository utilities
│       └── crud.go           # CRUD operations
├── usecases/                 # Business logic layer
│   ├── usecases.go           # Usecase container
│   ├── usecases_test.go      # Usecase tests
│   ├── arrows/               # Arrow business logic
│   │   ├── usecase.go        # Arrow usecase
│   │   └── usecase_test.go   # Arrow usecase tests
│   ├── quivers/              # Quiver business logic
│   │   ├── usecase.go        # Quiver usecase
│   │   └── usecase_test.go   # Quiver usecase tests
│   └── system/               # System business logic
│       ├── usecase.go        # System usecase
│       └── usecase_test.go   # System usecase tests
├── internal.go               # Dependency injection container
└── internal_test.go          # Internal tests
```

## Architecture Layers

### 1. **API Layer** (`/internal/api/`)
- **Purpose**: HTTP REST API endpoints for package management
- **Responsibilities**: Request/response handling, HTTP middleware
- **Dependencies**: Usecases layer
- **Key Files**: `api.go`, `v1/routes.go`, controllers

### 2. **Use Cases Layer** (`/internal/usecases/`)
- **Purpose**: Package management business logic and application rules
- **Responsibilities**: Orchestrating package operations
- **Dependencies**: Repository layer
- **Key Files**: `usecases.go`, domain-specific usecase files

### 3. **Repository Layer** (`/internal/repositories/`)
- **Purpose**: Package data access abstraction
- **Responsibilities**: Package persistence, repository integration
- **Dependencies**: Infrastructure layer
- **Key Files**: `repositories.go`, interface definitions

### 4. **Infrastructure Layer** (`/internal/infrastructure/`)
- **Purpose**: External concerns and technical implementation
- **Responsibilities**: Database connections, file system, external APIs
- **Dependencies**: None (bottom layer)
- **Key Files**: `infrastructure.go`, service implementations

### 5. **Models Layer** (`/internal/models/`)
- **Purpose**: Package management domain entities and business objects
- **Responsibilities**: Package data structures, business rules
- **Dependencies**: None (pure domain)
- **Key Files**: Domain entity files

### 6. **Core Layer** (`/internal/core/`)
- **Purpose**: Application core services
- **Responsibilities**: Configuration, logging, metadata
- **Dependencies**: None (foundational)
- **Key Files**: `core.go`, `config/`, `watcher/`

## File Naming Conventions

### Go Files
- **Main files**: `package_name.go`
- **Test files**: `package_name_test.go`
- **Interface files**: `interface.go`
- **Implementation files**: Descriptive names (e.g., `arrows.go`)

### Configuration Files
- **YAML configs**: `config.yaml`, `default.yaml`
- **Example configs**: `config.example.yaml`

### Documentation Files
- **All docs**: Located in `/docs/` directory
- **Naming**: Descriptive with hyphens (e.g., `development-setup.md`)

## Package Organization Principles

### 1. **Single Responsibility**
Each package has a single, well-defined responsibility:
- `models/arrow/` - Arrow package entities
- `models/quiver/` - Quiver server entities
- `repositories/arrows/` - Arrow data access
- `usecases/arrows/` - Arrow business logic

### 2. **Dependency Direction**
Dependencies flow inward following Clean Architecture:
```
API → Usecases → Repositories → Infrastructure
```

### 3. **Interface Segregation**
Small, focused interfaces:
- `ArrowsInterface` - Arrow data operations
- `QuiversInterface` - Quiver data operations
- `SystemInterface` - System data operations

### 4. **Dependency Injection**
All dependencies are injected through constructors:
```go
func NewUsecases(repositories *repositories.Repositories) *Usecases {
    return &Usecases{
        Arrows:  arrows.NewArrowsUsecase(repositories),
        Quivers: quivers.NewQuiversUsecase(repositories),
        System:  system.NewSystemUsecase(repositories),
    }
}
```

## Key Design Patterns

### 1. **Repository Pattern**
- Abstracts data access
- Enables testing with mocks
- Provides consistent data interface

### 2. **Dependency Injection**
- Centralized in `internal/internal.go`
- All dependencies injected through constructors
- Enables easy testing and configuration

### 3. **Clean Architecture**
- Clear separation of concerns
- Dependencies point inward
- Business logic isolated from technical details

### 4. **Command Pattern** (TUI)
- Commands parsed and handled consistently
- Event-driven architecture
- Separation of command parsing and execution

## Development Guidelines

### Adding New Features

1. **Models**: Add domain entities in `/internal/models/`
2. **Repositories**: Create data access in `/internal/repositories/`
3. **Use Cases**: Implement business logic in `/internal/usecases/`
4. **API**: Add endpoints in `/internal/api/v1/controllers/`
5. **Tests**: Add corresponding test files

### File Size Guidelines

- **Maximum file size**: 500 lines
- **If larger**: Split into multiple files within the same package
- **Consider**: Breaking into separate packages if functionality grows

### Import Organization

```go
// Standard library imports
import (
    "context"
    "fmt"
    "time"
)

// Third-party imports
import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// Internal imports
import (
    "github.com/rabbytesoftware/quiver/internal/models"
    "github.com/rabbytesoftware/quiver/internal/repositories"
)
```

## Testing Structure

### Test File Organization
- **Unit tests**: Co-located with source files (`*_test.go`)
- **Integration tests**: Separate test files for complex scenarios
- **Test data**: Use descriptive test data and helper functions

### Test Naming
```go
func TestNewArrowsUsecase(t *testing.T)           // Constructor tests
func TestArrowsUsecase_CreateArrow(t *testing.T) // Method tests
func TestArrowsUsecase_WithInvalidInput(t *testing.T) // Error case tests
```

## Configuration Management

### Configuration Files
- **Default config**: `internal/core/config/default.yaml`
- **Environment overrides**: Via environment variables
- **Runtime config**: Loaded through `config` package

### Configuration Structure
```yaml
config:
  api:
    host: 0.0.0.0
    port: 40257
  arrows:
    repositories: []
    install_dir: ./arrows
  watcher:
    enabled: true
    level: info
```

## Build and Development

### Makefile Commands
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linting
- `make clean` - Clean build artifacts

### Generated Files
- `/bin/` - Compiled binaries
- `/logs/` - Application logs
- `coverage.out` - Test coverage reports
- `coverage.html` - HTML coverage reports

## Best Practices

### 1. **Package Cohesion**
- Keep related functionality together
- Avoid circular dependencies
- Use clear package boundaries

### 2. **Error Handling**
- Always handle errors appropriately
- Use descriptive error messages
- Log errors at appropriate levels

### 3. **Documentation**
- Document public APIs
- Use meaningful variable names
- Add comments for complex logic

### 4. **Testing**
- Write tests for all public functions
- Aim for high test coverage
- Use descriptive test names

---

*For more details on specific components, see the [Architecture Overview](architecture-overview.md) and [Domain Models](domain-models.md) documentation.*
