# REST API Handlers Structure

This directory contains the restructured REST API handlers, organized by subsystem for better maintainability and modularity.

## Directory Structure

```
handlers/
├── README.md                  # This file
├── health/
│   └── health.go             # Health check endpoints
├── packages/
│   └── packages.go           # Legacy package management endpoints
├── arrows/
│   └── arrows.go             # Arrow package management endpoints
├── repositories/
│   └── repositories.go       # Repository management endpoints
└── server/
    └── server.go             # Server info/status endpoints
```

## Subsystems

### Health (`health/`)
- Health check endpoints
- Basic system status verification

### Packages (`packages/`)
- Legacy package management (for backward compatibility)
- Package listing, status, start/stop operations

### Arrows (`arrows/`)
- Arrow package management system
- Search, install, execute, uninstall, update, validate operations
- Status checking for installed arrows

### Repositories (`repositories/`)
- Repository management for arrow packages
- Add, remove, and list repositories

### Server (`server/`)
- Server information and status endpoints
- System metadata and operational status

## Usage

Each subsystem has its own handler struct that is instantiated in the main `Handlers` struct in `internal/server/handlers.go`. Routes are organized by subsystem in `internal/server/routes.go`.

## Benefits

1. **Better Organization**: Each subsystem has its own dedicated file and package
2. **Easier Maintenance**: Changes to one subsystem don't affect others
3. **Cleaner Imports**: Only import what you need for each handler
4. **Scalability**: Easy to add new subsystems or extend existing ones
5. **Testing**: Easier to write focused unit tests for each subsystem 