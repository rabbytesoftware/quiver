# Quiver Package Manager

Quiver now includes a comprehensive package manager that allows you to search, install, execute, and manage Arrow packages across multiple repositories.

## Features

- **Multi-repository support**: Add local directories or remote URLs as repositories
- **Dependency resolution**: Automatic installation and management of dependencies
- **Package lifecycle**: Install, execute, update, validate, and uninstall arrows
- **Local database**: Track installed packages and their status
- **REST API**: Full API support for all package operations
- **CLI interface**: Command-line tools for debugging and management

## Configuration

Update your `config.json` to include repository configuration:

```json
{
  "packages": {
    "repositories": [
      "./pkgs",
      "https://raw.githubusercontent.com/rabbytesoftware/quiver.pkgs/main"
    ],
    "directory": "./pkgs",
    "template_dir": "./template",
    "database_path": "./pkgs/packages.db"
  }
}
```

## REST API Endpoints

### Arrow Management

#### Search Arrows
```
GET /api/v1/arrows/search?q=query
```
Search for arrows across all repositories.

#### Install Arrow
```
POST /api/v1/arrows/{name}/install
Content-Type: application/json

{
  "variables": {
    "SERVER_HOSTNAME": "My CS2 Server",
    "MAX_PLAYERS": "16"
  }
}
```

#### Execute Arrow
```
POST /api/v1/arrows/{name}/execute
Content-Type: application/json

{
  "variables": {
    "SERVER_HOSTNAME": "My CS2 Server"
  }
}
```

#### Uninstall Arrow
```
DELETE /api/v1/arrows/{name}/uninstall
```

#### Update Arrow
```
PUT /api/v1/arrows/{name}/update
```

#### Validate Arrow
```
POST /api/v1/arrows/{name}/validate
```

#### Get Arrow Status
```
GET /api/v1/arrows/{name}/status
```

#### List Installed Arrows
```
GET /api/v1/arrows/installed
```

### Repository Management

#### List Repositories
```
GET /api/v1/repositories
```

#### Add Repository
```
POST /api/v1/repositories
Content-Type: application/json

{
  "repository": "https://github.com/example/arrows"
}
```

#### Remove Repository
```
DELETE /api/v1/repositories
Content-Type: application/json

{
  "repository": "https://github.com/example/arrows"
}
```

## CLI Commands

### Basic Package Operations

```bash
# Search for arrows
./quiver search minecraft

# Install an arrow
./quiver install cs2

# Execute an installed arrow
./quiver execute cs2

# Get arrow status
./quiver arrow-status cs2

# List installed arrows
./quiver installed

# Update an arrow
./quiver update cs2

# Validate an arrow
./quiver validate cs2

# Uninstall an arrow
./quiver uninstall cs2
```

### Repository Management

```bash
# List repositories
./quiver repo-list

# Add a repository
./quiver repo-add https://github.com/example/arrows

# Remove a repository
./quiver repo-remove https://github.com/example/arrows
```

### Server Management

```bash
# Check server health
./quiver health

# Get server info
./quiver info

# Get server status
./quiver server-status
```

## Arrow Structure

Arrows are defined in YAML files with the following structure:

```yaml
version: "0.1"

metadata:
  name: "Counter-Strike 2 SRCDS"
  description: "A basic Counter-Strike 2 Dedicated Server"
  mainteiners:
    - "https://char2cs.net"
  version: "0.0.1"

requirements:
  minimum:
    cpu_cores: 2
    ram_gb: 4
    disk_gb: 30
    network_mbps: 10

dependencies:
  - "steamcmd"

variables:
  - name: "SERVER_HOSTNAME"
    default: "CS2 Server hosted with Quiver"
  - name: "MAX_PLAYERS"
    default: 12
    min: 2
    max: 64

methods:
  install:
    linux:
      - "${steamcmd.execute} +login anonymous +force_install_dir ${INSTALL_DIR} +app_update 730 validate +quit"
  execute:
    linux:
      - "${INSTALL_DIR}/cs2 -dedicated -console +hostname ${SERVER_HOSTNAME} +maxplayers ${MAX_PLAYERS}"
  uninstall:
    linux:
      - "rm -rf ${INSTALL_DIR}"
  update:
    linux:
      - "${steamcmd.execute} +login anonymous +force_install_dir ${INSTALL_DIR} +app_update 730 validate +quit"
  validate:
    linux:
      - "test -f ${INSTALL_DIR}/cs2"
```

## Dependency Management

The package manager automatically handles dependencies:

1. **Installation**: Dependencies are installed before the main package
2. **Shared dependencies**: Common dependencies are only installed once
3. **Uninstallation**: Dependencies are only removed if no other packages depend on them
4. **Dependency tracking**: The system tracks which packages depend on others

## Variable System

Variables can be defined in arrows and customized during installation:

- **Default values**: Specified in the arrow definition
- **User overrides**: Provided during installation or execution
- **Sensitive variables**: Marked as sensitive (passwords, API keys)
- **Type validation**: Support for strings, integers, and booleans
- **Value constraints**: Min/max values and predefined options

### System Variables

Quiver automatically provides system variables that arrows can reference:

- **`INSTALL_DIR`**: The directory where the arrow's files are installed (e.g., `./pkgs/chat/` for `chat.yaml`)
- **`INSTALL_PATH`**: Deprecated alias for `INSTALL_DIR` (maintained for backward compatibility)

## Package States

Packages can be in the following states:

- `installed`: Package is installed but not running
- `running`: Package is currently executing
- `stopped`: Package was running but has stopped
- `error`: Package encountered an error during execution

## Repository Types

### Local Repositories
- Local directory containing arrow YAML files
- Example: `./pkgs`, `/path/to/arrows`

### Remote Repositories
- HTTP/HTTPS URLs to repositories
- Example: `https://raw.githubusercontent.com/rabbytesoftware/quiver.pkgs/main`

## Error Handling

The package manager provides detailed error messages for:

- Missing dependencies
- Installation failures
- Execution errors
- Network issues
- Configuration problems
- Permission errors

## Best Practices

1. **Repository Organization**: Keep arrows organized in clear directory structures
2. **Dependency Management**: Minimize dependencies and use common base dependencies
3. **Variable Naming**: Use clear, descriptive variable names
4. **Error Handling**: Include proper error handling in arrow methods
5. **Documentation**: Document arrow requirements and usage
6. **Testing**: Validate arrows thoroughly before distribution

## Troubleshooting

### Common Issues

1. **Installation fails**: Check dependencies and system requirements
2. **Execution fails**: Verify variables and installation paths
3. **Repository access**: Ensure network connectivity for remote repositories
4. **Permission errors**: Check file system permissions for installation directory

### Debug Commands

```bash
# Get detailed server status
./quiver server-status

# Check arrow installation status
./quiver arrow-status arrow-name

# Validate arrow installation
./quiver validate arrow-name
```

## Examples

### Installing and Running CS2 Server

```bash
# Search for CS2 arrow
./quiver search cs2

# Install CS2 with custom variables
curl -X POST http://localhost:45565/api/v1/arrows/cs2/install \
  -H "Content-Type: application/json" \
  -d '{
    "variables": {
      "SERVER_HOSTNAME": "My Awesome CS2 Server",
      "MAX_PLAYERS": "20",
      "DEFAULT_MAP": "de_mirage"
    }
  }'

# Execute the server
./quiver execute cs2

# Check status
./quiver arrow-status cs2
```

### Managing Repositories

```bash
# Add a new repository
./quiver repo-add https://github.com/example/game-arrows

# List all repositories
./quiver repo-list

# Search for arrows in all repositories
./quiver search minecraft
```

This comprehensive package manager system provides a robust foundation for managing game server and application deployments with Quiver. 